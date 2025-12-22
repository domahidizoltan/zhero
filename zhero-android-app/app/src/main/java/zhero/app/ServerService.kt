package zhero.app

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.os.Environment
import android.os.IBinder
import android.util.Log
import androidx.core.app.NotificationCompat
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.Job
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import server.Server
import server.Server_
import java.io.ByteArrayOutputStream
import java.io.InputStreamReader
import java.io.OutputStream
import java.io.PrintStream
import java.lang.StringBuilder

class ServerService : Service() {

    private var serverInstance: Server_? = null
    private val serviceJob = Job()
    private val serviceScope = CoroutineScope(Dispatchers.IO + serviceJob)

    // For capturing Go's stderr
    private lateinit var goStderrStream: ByteArrayOutputStream
    private lateinit var goStderrPrintStream: PrintStream
    private var stderrMonitorJob: Job? = null

    companion object {
        const val CHANNEL_ID = "ServerServiceChannel"
        const val NOTIFICATION_ID = 1
        const val TAG = "ServerService" // This TAG is used for filtering logs in MainActivity
    }

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
        val notification = createNotification()
        startForeground(NOTIFICATION_ID, notification)

        // Initialize stderr capture
        goStderrStream = ByteArrayOutputStream()
        goStderrPrintStream = PrintStream(goStderrStream)

        initializeServerInstance()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        when (intent?.action) {
            "START_SERVER" -> startServer()
            "STOP_SERVER" -> stopSelf()
        }
        return START_NOT_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        stopServer()
        serviceJob.cancel()
        // Clean up stderr print stream
        goStderrPrintStream.close()
        stderrMonitorJob?.cancel()
    }

    override fun onBind(intent: Intent?): IBinder? {
        // This service does not allow binding
        return null
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val serviceChannel = NotificationChannel(
                CHANNEL_ID,
                "Zhero Server Channel",
                NotificationManager.IMPORTANCE_DEFAULT
            )
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(serviceChannel)
        }
    }

    private fun createNotification(): Notification {
        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle("Zhero Server Running")
            .setContentText("Your local Zhero server is active in the background.")
            .setSmallIcon(R.mipmap.ic_launcher)
            .setPriority(NotificationCompat.PRIORITY_LOW)
            .build()
    }

    private fun initializeServerInstance() {
        Log.d(TAG, "Attempting to initialize server instance in service...")

        serviceScope.launch {
            try {
                serverInstance = Server.new_()
                // Set the stderr for the Go runtime to our PrintStream
//                serverInstance?.setStderr(goStderrPrintStream)

                // Ensure the path ends with a slash if it's a directory
                val basePath = Environment.getExternalStorageDirectory().path + "/"
                serverInstance?.setAbsolutePath(basePath)

                // Start monitoring Go's stderr
                startStderrMonitoring()

            } catch (e: Exception) {
                Log.e(TAG, "Error initializing server instance: ${e.message}", e)
            }
        }
    }

    private fun startStderrMonitoring() {
        if (stderrMonitorJob != null && stderrMonitorJob!!.isActive) {
            return // Already monitoring
        }
        stderrMonitorJob = serviceScope.launch {
            while (true) {
                try {
                    // Check if there's anything to read from the stream
                    if (goStderrStream.size() > 0) {
                        val output = goStderrStream.toString()
                        goStderrStream.reset() // Clear the buffer after reading
                        if (output.isNotBlank()) {
                            Log.e(TAG, "Go stderr: $output")
                        }
                    }
                } catch (e: Exception) {
                    Log.e(TAG, "Error monitoring Go stderr: ${e.message}", e)
                }
                // Wait a bit before checking again to avoid busy-waiting
                delay(100)
            }
        }
    }

    private fun startServer() {
        if (serverInstance == null) {
            Log.w(TAG, "Server instance is null, attempting to re-initialize before starting.")
            initializeServerInstance()
            // Give a moment for initialization to complete before trying to start
            serviceScope.launch {
                delay(500) // Small delay
                if (serverInstance != null) {
                    startServerInternal()
                } else {
                    Log.e(TAG, "Server instance still null after re-initialization attempt.")
                }
            }
        } else {
            startServerInternal()
        }
    }

    private fun startServerInternal() {
        serviceScope.launch {
            try {
                serverInstance?.start()
            } catch (e: Exception) {
                Log.e(TAG, "Error starting server: ${e.message}", e)
            }
        }
    }

    private fun stopServer() {
        if (serverInstance == null) {
            Log.w(TAG, "Server instance is null, no need to stop.")
            return
        }
        serviceScope.launch {
            try {
                serverInstance?.stop()
            } catch (e: Exception) {
                Log.e(TAG, "Error stopping server: ${e.message}", e)
            } finally {
                serverInstance = null
            }
        }
    }
}
