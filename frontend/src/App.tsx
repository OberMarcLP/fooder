import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom';
import { Home, Tag, Utensils } from 'lucide-react';
import { useTheme } from './hooks/useTheme';
import { ThemeToggle } from './components/ThemeToggle';
import { GlobalSearch } from './components/GlobalSearch';
import { HomePage } from './pages/HomePage';
import { CategoriesPage } from './pages/CategoriesPage';
import { FoodTypesPage } from './pages/FoodTypesPage';
import { Category, FoodType, RestaurantFilters, getCategories, getFoodTypes } from './services/api';

function App() {
  const { isDark, toggleTheme } = useTheme();
  const [categories, setCategories] = useState<Category[]>([]);
  const [foodTypes, setFoodTypes] = useState<FoodType[]>([]);
  const [filters, setFilters] = useState<RestaurantFilters>({});

  useEffect(() => {
    const loadFiltersData = async () => {
      try {
        const [cats, fts] = await Promise.all([getCategories(), getFoodTypes()]);
        setCategories(cats);
        setFoodTypes(fts);
      } catch (error) {
      }
    };
    loadFiltersData();
  }, []);

  return (
    <BrowserRouter>
      <div className="min-h-screen">
        <nav className="bg-white dark:bg-gray-800 shadow-lg sticky top-0 z-40">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex flex-col gap-3 py-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-8">
                  <span className="text-xl font-bold text-blue-600 dark:text-blue-400">
                    The Nom Database
                  </span>
                  <div className="hidden lg:flex gap-1">
                    <NavLink
                      to="/"
                      className={({ isActive }) =>
                        `flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
                          isActive
                            ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                            : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                        }`
                      }
                    >
                      <Home className="w-4 h-4" />
                      <span>Restaurants</span>
                    </NavLink>
                    <NavLink
                      to="/categories"
                      className={({ isActive }) =>
                        `flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
                          isActive
                            ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                            : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                        }`
                      }
                    >
                      <Tag className="w-4 h-4" />
                      <span>Categories</span>
                    </NavLink>
                    <NavLink
                      to="/food-types"
                      className={({ isActive }) =>
                        `flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
                          isActive
                            ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                            : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                        }`
                      }
                    >
                      <Utensils className="w-4 h-4" />
                      <span>Food Types</span>
                    </NavLink>
                  </div>
                </div>
                <ThemeToggle isDark={isDark} onToggle={toggleTheme} />
              </div>

              {/* Global Search Bar */}
              <div className="w-full">
                <GlobalSearch
                  categories={categories}
                  foodTypes={foodTypes}
                  filters={filters}
                  onFiltersChange={setFilters}
                />
              </div>

              {/* Mobile Navigation */}
              <div className="flex lg:hidden gap-1 overflow-x-auto pb-1">
                <NavLink
                  to="/"
                  className={({ isActive }) =>
                    `flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap ${
                      isActive
                        ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                        : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                    }`
                  }
                >
                  <Home className="w-4 h-4" />
                  <span className="text-sm">Restaurants</span>
                </NavLink>
                <NavLink
                  to="/categories"
                  className={({ isActive }) =>
                    `flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap ${
                      isActive
                        ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                        : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                    }`
                  }
                >
                  <Tag className="w-4 h-4" />
                  <span className="text-sm">Categories</span>
                </NavLink>
                <NavLink
                  to="/food-types"
                  className={({ isActive }) =>
                    `flex items-center gap-2 px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap ${
                      isActive
                        ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                        : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                    }`
                  }
                >
                  <Utensils className="w-4 h-4" />
                  <span className="text-sm">Food Types</span>
                </NavLink>
              </div>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route path="/" element={<HomePage filters={filters} />} />
            <Route path="/categories" element={<CategoriesPage />} />
            <Route path="/food-types" element={<FoodTypesPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
