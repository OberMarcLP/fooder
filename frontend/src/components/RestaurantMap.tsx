import { GoogleMap, Marker, useJsApiLoader } from '@react-google-maps/api';

interface RestaurantMapProps {
  latitude: number;
  longitude: number;
  name: string;
}

const containerStyle = {
  width: '100%',
  height: '300px',
  borderRadius: '0.5rem',
};

export function RestaurantMap({ latitude, longitude, name }: RestaurantMapProps) {
  const { isLoaded, loadError } = useJsApiLoader({
    googleMapsApiKey: import.meta.env.VITE_GOOGLE_MAPS_API_KEY || '',
  });

  const center = { lat: latitude, lng: longitude };

  if (loadError) {
    return (
      <div className="w-full h-[300px] bg-gray-200 dark:bg-gray-700 rounded-lg flex items-center justify-center">
        <p className="text-gray-500 dark:text-gray-400">Failed to load map</p>
      </div>
    );
  }

  if (!isLoaded) {
    return (
      <div className="w-full h-[300px] bg-gray-200 dark:bg-gray-700 rounded-lg flex items-center justify-center animate-pulse">
        <p className="text-gray-500 dark:text-gray-400">Loading map...</p>
      </div>
    );
  }

  return (
    <div className="relative">
      <GoogleMap mapContainerStyle={containerStyle} center={center} zoom={15}>
        <Marker position={center} title={name} />
      </GoogleMap>
      <a
        href={`https://www.google.com/maps/dir/?api=1&destination=${latitude},${longitude}`}
        target="_blank"
        rel="noopener noreferrer"
        className="absolute bottom-4 right-4 btn btn-primary text-sm"
      >
        Get Directions
      </a>
    </div>
  );
}
