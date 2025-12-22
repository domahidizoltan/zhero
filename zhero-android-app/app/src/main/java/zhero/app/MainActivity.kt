package zhero.app

import android.Manifest
import android.app.ActivityManager
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.net.Uri
import android.os.Build
import android.os.Bundle
import android.os.Environment
import android.provider.Settings
import android.util.Log
import android.widget.Button
import android.widget.ScrollView
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import com.google.android.material.snackbar.Snackbar
import zhero.app.databinding.ActivityMainBinding
import java.io.BufferedReader
import java.io.InputStreamReader
import java.lang.StringBuilder

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var logTextView: TextView
    private lateinit var logScrollView: ScrollView
    private lateinit var serverToggleButton: Button

    private val PERMISSION_REQUEST_CODE = 100
    private var isServerRunning = false

    // Tag for filtering log messages. Changed to "GoLog" as observed.
    private val GO_LOG_TAG = "GoLog"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        logTextView = binding.logTextView
        logScrollView = binding.logScrollView
        serverToggleButton = binding.serverToggleButton

        serverToggleButton.setOnClickListener {
            if (isServerRunning) {
                stopServerService()
            } else {
                startServerService()
            }
        }

        checkAndRequestPermissions()
        // Start monitoring logs after permissions are handled and UI is set up
        startLogMonitoring()
    }

    override fun onResume() {
        super.onResume()
        // Ensure UI is updated when activity resumes
        updateServerButtonUI()
        // Re-check permissions if necessary
        checkAndRequestPermissions()
    }

    private fun startServerService() {
        Log.d("MainActivity", "Attempting to start ServerService...")
        serverToggleButton.isEnabled = false
        val serviceIntent = Intent(this, ServerService::class.java).apply {
            action = "START_SERVER"
        }
        ContextCompat.startForegroundService(this, serviceIntent)
        isServerRunning = true
        updateServerButtonUI()
        // Re-enable button after a short delay to prevent rapid toggling
        serverToggleButton.postDelayed({ serverToggleButton.isEnabled = true }, 1000)
    }

    private fun stopServerService() {
        Log.d("MainActivity", "Attempting to stop ServerService...")
        serverToggleButton.isEnabled = false
        val serviceIntent = Intent(this, ServerService::class.java).apply {
            action = "STOP_SERVER"
        }
        stopService(serviceIntent)
        isServerRunning = false
        updateServerButtonUI()
        // Re-enable button after a short delay
        serverToggleButton.postDelayed({ serverToggleButton.isEnabled = true }, 1000)
    }

    private fun updateServerButtonUI() {
        if (isServerRunning) {
            serverToggleButton.setText(R.string.stop_server)
            serverToggleButton.setBackgroundColor(ContextCompat.getColor(this, R.color.red_button))
        } else {
            serverToggleButton.setText(R.string.start_server)
            serverToggleButton.setBackgroundColor(ContextCompat.getColor(this, R.color.green_button))
        }
    }

    private fun isServiceRunning(serviceClass: Class<*>): Boolean {
        val manager = getSystemService(Context.ACTIVITY_SERVICE) as ActivityManager
        for (service in manager.getRunningServices(Integer.MAX_VALUE)) {
            if (serviceClass.name == service.service.className) {
                return true
            }
        }
        return false
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)
        if (requestCode == PERMISSION_REQUEST_CODE) {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
                if (Environment.isExternalStorageManager()) {
                    Snackbar.make(binding.root, "All Files Access permission granted", Snackbar.LENGTH_SHORT).show()
                } else {
                    Snackbar.make(binding.root, "All Files Access permission denied. Please enable it in Settings.", Snackbar.LENGTH_LONG)
                        .setAction("Settings") {
                            val intent = Intent(Settings.ACTION_MANAGE_APP_ALL_FILES_ACCESS_PERMISSION)
                            startActivity(intent)
                        }.show()
                }
            }
        }
    }

    private fun checkAndRequestPermissions() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) { // API 30+
            if (!Environment.isExternalStorageManager()) {
                val intent = Intent(Settings.ACTION_MANAGE_APP_ALL_FILES_ACCESS_PERMISSION, Uri.parse("package:$packageName"))
                startActivityForResult(intent, PERMISSION_REQUEST_CODE)
            }
        } else { // < API 30
            if (ContextCompat.checkSelfPermission(this, Manifest.permission.WRITE_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED) {
                ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.WRITE_EXTERNAL_STORAGE), PERMISSION_REQUEST_CODE)
            }
        }
    }

    override fun onRequestPermissionsResult(requestCode: Int, permissions: Array<out String>, grantResults: IntArray) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == PERMISSION_REQUEST_CODE) {
            if (Build.VERSION.SDK_INT < Build.VERSION_CODES.R) { // Only for API < 30
                if (grantResults.isNotEmpty() && grantResults[0] == PackageManager.PERMISSION_GRANTED) {
                    Snackbar.make(binding.root, "Write External Storage permission granted", Snackbar.LENGTH_SHORT).show()
                } else {
                    Snackbar.make(binding.root, "Write External Storage permission denied. Please enable it in Settings.", Snackbar.LENGTH_LONG)
                        .setAction("Settings") {
                            val intent = Intent(Settings.ACTION_APPLICATION_DETAILS_SETTINGS)
                            val uri = Uri.fromParts("package", packageName, null)
                            intent.data = uri
                            startActivity(intent)
                        }.show()
                }
            }
        }
    }

    // --- Log Monitoring ---
    private fun startLogMonitoring() {
        // This function will set up a mechanism to read from Logcat and append to logTextView.
        // We'll use a background thread to continuously poll Logcat for messages tagged with GO_LOG_TAG.

        Thread {
            try {
                // Execute logcat command to capture messages with the specified tag.
                // '-s' silences the tag itself, so we only get the message content.
                // This command will run continuously as long as the app is in the foreground and this thread is alive.
                val process = Runtime.getRuntime().exec("logcat -s $GO_LOG_TAG")
                val reader = BufferedReader(InputStreamReader(process.inputStream))
                val logBuffer = StringBuilder()
                var line: String?

                while (true) {
                    line = reader.readLine()
                    if (line != null) {
                        // Append the log line to our buffer
                        logBuffer.append(line).append("\n")
                        // Update the UI on the main thread
                        runOnUiThread {
                            logTextView.text = logBuffer.toString()
                            // Auto-scroll to the bottom
                            logScrollView.post {
                                logScrollView.fullScroll(ScrollView.FOCUS_DOWN)
                            }
                        }
                    } else {
                        // If readLine() returns null, the logcat process might have terminated.
                        // This could happen if the app is backgrounded or if logcat itself stops.
                        // We can add a small delay and retry or break. For now, let's break.
                        break
                    }
                }
                // Ensure the process is cleaned up if the loop breaks
                process.destroy()
            } catch (e: Exception) {
                Log.e("MainActivity", "Error reading logcat: ${e.message}", e)
                runOnUiThread {
                    logTextView.append("\nError reading logs: ${e.message}")
                }
            }
        }.start()
    }
}
