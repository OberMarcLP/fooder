import { useState } from 'react';
import { Search, Loader2 } from 'lucide-react';
import { searchPlaces, GooglePlaceResult } from '../services/api';

interface PlaceSearchProps {
  onSelect: (place: GooglePlaceResult) => void;
}

export function PlaceSearch({ onSelect }: PlaceSearchProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<GooglePlaceResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    if (!query.trim()) return;

    setLoading(true);
    setError(null);
    try {
      const places = await searchPlaces(query);
      setResults(places);
      if (places.length === 0) {
        setError('No restaurants found. Try a different search.');
      }
    } catch (err) {
      setError('Failed to search. Please try again.');
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSelect = (place: GooglePlaceResult) => {
    onSelect(place);
    setResults([]);
    setQuery('');
  };

  return (
    <div className="relative">
      <div className="flex gap-2">
        <div className="relative flex-1">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
            placeholder="Search restaurant name..."
            className="input !pl-12"
          />
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        </div>
        <button
          onClick={handleSearch}
          disabled={loading || !query.trim()}
          className="btn btn-primary flex items-center gap-2"
        >
          {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Search'}
        </button>
      </div>

      {error && <p className="text-red-500 text-sm mt-2">{error}</p>}

      {results.length > 0 && (
        <div className="absolute z-10 w-full mt-2 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 max-h-80 overflow-y-auto">
          {results.map((place) => (
            <button
              key={place.place_id}
              onClick={() => handleSelect(place)}
              className="w-full text-left px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-700 border-b border-gray-100 dark:border-gray-700 last:border-0"
            >
              <p className="font-medium">{place.name}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">{place.address}</p>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
