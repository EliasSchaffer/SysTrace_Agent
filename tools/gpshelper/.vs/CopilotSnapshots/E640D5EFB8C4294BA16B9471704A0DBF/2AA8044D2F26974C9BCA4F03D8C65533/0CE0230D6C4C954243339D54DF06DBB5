using System;
using System.Diagnostics;
using System.Text.Json;
using System.Threading.Tasks;
using Windows.Devices.Geolocation;

class Program
{
    static async Task Main()
    {
        var geolocator = new Geolocator();

        double latitude = 0, longitude = 0, altitude = 0, accuracy = 0;

        try
        {
            var accessStatus = await Geolocator.RequestAccessAsync();

            if (accessStatus != GeolocationAccessStatus.Allowed)
            {
                Console.Error.WriteLine($"Location access status: {accessStatus}");
                Console.Error.WriteLine("\nTo enable location access:");
                Console.Error.WriteLine("1. Open Windows Settings (Win + I)");
                Console.Error.WriteLine("2. Go to Privacy & Security > Location");
                Console.Error.WriteLine("3. Turn ON 'Location services'");
                Console.Error.WriteLine("4. Turn ON 'Let apps access your location'");
                Console.Error.WriteLine("5. Scroll down and enable for 'gpshelper'");
                Console.Error.WriteLine("\nOpening Windows Settings now...");

                try
                {
                    Process.Start(new ProcessStartInfo
                    {
                        FileName = "ms-settings:privacy-location",
                        UseShellExecute = true
                    });
                }
                catch { }
            }
            else
            {
                Console.Error.WriteLine("Getting location...");
                var position = await geolocator.GetGeopositionAsync(
                    maximumAge: TimeSpan.FromMinutes(5),
                    timeout: TimeSpan.FromSeconds(10));

                latitude = position.Coordinate.Point.Position.Latitude;
                longitude = position.Coordinate.Point.Position.Longitude;
                altitude = position.Coordinate.Point.Position.Altitude;
                accuracy = position.Coordinate.Accuracy;

                Console.Error.WriteLine("Location retrieved successfully!");
            }
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"Error getting location: {ex.Message}");
            Console.Error.WriteLine($"Exception type: {ex.GetType().Name}");
        }

        var gps = new
        {
            latitude = latitude,
            longitude = longitude,
            altitude = altitude,
            accuracy = accuracy
        };

        Console.WriteLine(JsonSerializer.Serialize(gps));
    }
}