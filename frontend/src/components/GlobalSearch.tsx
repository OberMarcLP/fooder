import { useState, useEffect, useRef } from 'react';
import { Search, X, Lightbulb, Filter } from 'lucide-react';
import { globalSearch, Restaurant, Category, FoodType, RestaurantFilters } from '../services/api';
import { useNavigate } from 'react-router-dom';
import { SearchFilters } from './SearchFilters';

interface GlobalSearchProps {
  categories: Category[];
  foodTypes: FoodType[];
  filters: RestaurantFilters;
  onFiltersChange: (filters: RestaurantFilters) => void;
}

export function GlobalSearch({ categories, foodTypes, filters, onFiltersChange }: GlobalSearchProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Restaurant[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [showFilters, setShowFilters] = useState(false);
  const searchRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    const searchDebounce = setTimeout(async () => {
      if (query.trim().length >= 2) {
        setIsLoading(true);
        try {
          const data = await globalSearch(query);
          setResults(data);
          setIsOpen(true);
        } catch (error) {
          setResults([]);
        } finally {
          setIsLoading(false);
        }
      } else {
        setResults([]);
        setIsOpen(false);
      }
    }, 300);

    return () => clearTimeout(searchDebounce);
  }, [query]);

  const handleResultClick = () => {
    // Navigate to home page - suggestions will be shown there too
    navigate(`/`);
    setIsOpen(false);
    setQuery('');
  };

  const handleClear = () => {
    setQuery('');
    setResults([]);
    setIsOpen(false);
  };

  const hasActiveFilters = filters.category_id || (filters.food_type_ids && filters.food_type_ids.length > 0) || filters.radius;

  return (
    <div ref={searchRef} className="relative w-full">
      <div className="flex gap-2 items-center mb-2">
        <div className="relative flex-1">
          <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 transition-colors peer-focus:text-blue-500 w-5 h-5" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onFocus={() => query.length >= 2 && setIsOpen(true)}
            placeholder="Search restaurants and suggestions..."
            className="input-glass w-full pl-11 pr-11 peer focus:shadow-lg focus:shadow-blue-500/20"
          />
          {query && (
            <button
              onClick={handleClear}
              className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              <X className="w-5 h-5" />
            </button>
          )}
        </div>
        <button
          onClick={() => setShowFilters(!showFilters)}
          className={`relative p-2 rounded-xl transition-all duration-300 ${
            showFilters
              ? 'bg-gradient-to-r from-blue-500/20 to-purple-500/20 backdrop-blur-md border border-blue-500/30 shadow-lg shadow-blue-500/20'
              : 'btn-glass'
          }`}
        >
          <Filter className="w-5 h-5" />
          {hasActiveFilters && (
            <span className="absolute top-0 right-0 w-2 h-2 bg-blue-500 rounded-full animate-pulse" />
          )}
        </button>
      </div>

      {showFilters && (
        <SearchFilters
          categories={categories}
          foodTypes={foodTypes}
          filters={filters}
          onFiltersChange={onFiltersChange}
        />
      )}

      {isOpen && (
        <div className="absolute z-50 w-full mt-3 bg-white/95 dark:bg-gray-800/95 backdrop-blur-2xl border border-white/30 dark:border-white/10 rounded-2xl shadow-2xl shadow-black/20 max-h-96 overflow-y-auto animate-slide-down">
          {isLoading ? (
            <div className="p-4 text-center text-gray-500 dark:text-gray-400">
              Searching...
            </div>
          ) : results.length > 0 ? (
            <ul className="divide-y divide-gray-200 dark:divide-gray-700">
              {results.map((restaurant, index) => (
                <li
                  key={`${restaurant.is_suggestion ? 's' : 'r'}-${restaurant.is_suggestion ? restaurant.suggestion_id : restaurant.id}-${index}`}
                  onClick={() => handleResultClick()}
                  className="p-4 hover:bg-white/40 dark:hover:bg-white/10 cursor-pointer transition-all duration-200 hover:scale-[1.01]"
                >
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h3 className="font-semibold text-gray-900 dark:text-white truncate">
                          {restaurant.name}
                        </h3>
                        {restaurant.is_suggestion && (
                          <span className="badge-suggestion flex-shrink-0">
                            <Lightbulb className="w-3 h-3" />
                            Suggestion
                          </span>
                        )}
                      </div>
                      {restaurant.address && (
                        <p className="text-sm text-gray-600 dark:text-gray-400 truncate mt-1">
                          {restaurant.address}
                        </p>
                      )}
                      {restaurant.category && (
                        <p className="text-xs text-gray-500 dark:text-gray-500 mt-1">
                          {restaurant.category.name}
                        </p>
                      )}
                    </div>
                    {!restaurant.is_suggestion && restaurant.avg_rating && (
                      <div className="text-right flex-shrink-0">
                        <div className="text-sm font-semibold text-gray-900 dark:text-white">
                          ★ {restaurant.avg_rating.overall.toFixed(1)}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          {restaurant.avg_rating.count} reviews
                        </div>
                      </div>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          ) : query.length >= 2 ? (
            <div className="p-4 text-center text-gray-500 dark:text-gray-400">
              No results found for "{query}"
            </div>
          ) : null}
        </div>
      )}
    </div>
  );
}
