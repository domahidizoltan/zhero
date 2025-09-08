document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');
    const resultDisplay = document.getElementById('result-display');
    const searchForm = document.getElementById('search-form');

    // Sample data for autocomplete
    const data = [
        'Thing', 'Intangible', 'Action', 'Event', 'Product', 'Service', 'Place',
        'Organization', 'Person', 'CreativeWork', 'MedicalEntity',
        'Accommodation', 'AdministrativeArea', 'Airport', 'AmusementPark',
        'Apartment', 'Aquarium', 'ArtGallery', 'Attorney', 'Bakery',
        'BankAccount', 'BarOrPub', 'Beach', 'BeautySalon', 'BedAndBreakfast',
        'BikeStore', 'BookStore', 'BowlingAlley', 'Brewery', 'Bridge',
        'BuddhistTemple', 'BusStation', 'BusStop', 'CafeOrCoffeeShop',
        'Campground', 'Canal', 'Casino', 'Cemetery', 'Church', 'City',
        'CityHall', 'CivicStructure', 'ClothingStore', 'ComedyClub',
        'ComputerStore', 'Continent', 'ConvenienceStore', 'Country',
        'Courthouse', 'Crematorium', 'DaySpa', 'Dentist', 'DepartmentStore',
        'Distillery', 'DryCleaningOrLaundry', 'Electrician', 'ElectronicsStore',
        'Embassy', 'EmergencyService', 'EmploymentAgency', 'EntertainmentBusiness',
        'EventVenue', 'ExerciseGym', 'Factory', 'FastFoodRestaurant',
        'FinancialService', 'FireStation', 'Florist', 'FoodEstablishment',
        'FuneralHome', 'FurnitureStore', 'GardenStore', 'GasStation',
        'GeneralContractor', 'GolfCourse', 'GovernmentBuilding', 'GovernmentOffice',
        'GroceryStore', 'HairSalon', 'HardwareStore', 'HealthAndBeautyBusiness',
        'HealthClub', 'HinduTemple', 'HobbyShop', 'HomeAndConstructionBusiness',
        'HomeGoodsStore', 'Hospital', 'Hostel', 'Hotel', 'HousePainter',
        'IceCreamShop', 'InsuranceAgency', 'InternetCafe', 'JewelryStore',
        'LakeBodyOfWater', 'Landform', 'LandmarksOrHistoricalBuildings', 'Library',
        'LiquorStore', 'LocalBusiness', 'Locksmith', 'LodgingBusiness',
        'MedicalBusiness', 'MedicalClinic', 'MensClothingStore', 'MobilePhoneStore',
        'Mosque', 'Motel', 'MotorcycleDealer', 'MotorcycleRepair', 'MovieRentalStore',
        'MovieTheater', 'MovingCompany', 'Museum', 'MusicStore', 'MusicVenue',
        'NightClub', 'Notary', 'NursingHome', 'OceanBodyOfWater', 'OfficeEquipmentStore',
        'OutletStore', 'PaintingService', 'Park', 'ParkingFacility', 'PawnShop',
        'PerformingArtsTheater', 'PetStore', 'Pharmacy', 'Physician',
        'PlaceOfWorship', 'Playground', 'Plumber', 'PoliceStation', 'Pond',
        'PostOffice', 'ProfessionalService', 'PublicSwimmingPool', 'RadioStation',
        'RealEstateAgent', 'RecyclingCenter', 'Reservoir', 'Restaurant',
        'RiverBodyOfWater', 'RoofingContractor', 'School', 'SeaBodyOfWater',
        'SelfStorage', 'ShoeStore', 'ShoppingCenter', 'SkiResort', 'SportingGoodsStore',
        'SportsActivityLocation', 'SportsClub', 'StadiumOrArena', 'State',
        'Store', 'SubwayStation', 'Synagogue', 'TattooParlor', 'TaxiStand',
        'TelevisionStation', 'TennisComplex', 'TheaterGroup', 'TireShop',
        'TouristAttraction', 'TouristInformationCenter', 'Tower', 'TownSquare',
        'TrainStation', 'TravelAgency', 'University', 'VeterinaryCare',
        'Village', 'Volcano', 'Waterfall', 'WholesaleStore', 'Winery', 'Zoo'
    ];

    function filterData(query) {
        return data.filter(item => item.toLowerCase().includes(query.toLowerCase()));
    }

    function showResults(results) {
        searchResults.innerHTML = '';
        if (results.length > 0) {
            results.forEach(item => {
                const li = document.createElement('li');
                li.className = 'p-2 hover:bg-base-200 cursor-pointer';
                li.textContent = item;
                li.addEventListener('click', () => {
                    searchInput.value = item;
                    searchResults.classList.add('hidden');
                    resultDisplay.textContent = `You selected: ${item}`;
                });
                searchResults.appendChild(li);
            });
            searchResults.classList.remove('hidden');
        } else {
            searchResults.classList.add('hidden');
        }
    }

    searchInput.addEventListener('input', () => {
        const query = searchInput.value;
        if (query.length > 0) {
            const filteredData = filterData(query);
            showResults(filteredData);
        } else {
            searchResults.classList.add('hidden');
        }
    });

    // Hide results when clicking outside
    document.addEventListener('click', (e) => {
        if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
            searchResults.classList.add('hidden');
        }
    });

    searchForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const value = searchInput.value;
        if (value) {
            resultDisplay.textContent = `You submitted: ${value}`;
        }
        searchResults.classList.add('hidden');
    });
});
