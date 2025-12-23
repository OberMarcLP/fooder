import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom';
import { Home, Settings } from 'lucide-react';
import { useTheme } from './hooks/useTheme';
import { ThemeToggle } from './components/ThemeToggle';
import { GlobalSearch } from './components/GlobalSearch';
import { HomePage } from './pages/HomePage';
import { SettingsPage } from './pages/SettingsPage';
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
        <nav className="nav-glass sticky top-0 z-40">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex flex-col gap-3 py-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-8">
                  <span className="text-xl font-bold text-gradient">
                    The Nom Database
                  </span>
                  <div className="hidden lg:flex gap-2">
                    <NavLink
                      to="/"
                      className={({ isActive }) =>
                        `flex items-center gap-2 px-4 py-2 rounded-full transition-all duration-300 ${
                          isActive
                            ? 'bg-gradient-to-r from-blue-500/20 to-purple-500/20 backdrop-blur-md border border-blue-500/30 shadow-lg shadow-blue-500/20'
                            : 'hover:bg-white/20 dark:hover:bg-white/10 hover:backdrop-blur-md hover:shadow-md'
                        }`
                      }
                    >
                      <Home className="w-4 h-4" />
                      <span>Restaurants</span>
                    </NavLink>
                    <NavLink
                      to="/settings"
                      className={({ isActive }) =>
                        `flex items-center gap-2 px-4 py-2 rounded-full transition-all duration-300 ${
                          isActive
                            ? 'bg-gradient-to-r from-blue-500/20 to-purple-500/20 backdrop-blur-md border border-blue-500/30 shadow-lg shadow-blue-500/20'
                            : 'hover:bg-white/20 dark:hover:bg-white/10 hover:backdrop-blur-md hover:shadow-md'
                        }`
                      }
                    >
                      <Settings className="w-4 h-4" />
                      <span>Settings</span>
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
              <div className="flex lg:hidden gap-2 overflow-x-auto pb-1">
                <NavLink
                  to="/"
                  className={({ isActive }) =>
                    `flex items-center gap-2 px-3 py-1.5 rounded-full transition-all duration-300 whitespace-nowrap ${
                      isActive
                        ? 'bg-gradient-to-r from-blue-500/20 to-purple-500/20 backdrop-blur-md border border-blue-500/30 shadow-lg shadow-blue-500/20'
                        : 'hover:bg-white/20 dark:hover:bg-white/10 hover:backdrop-blur-md hover:shadow-md'
                    }`
                  }
                >
                  <Home className="w-4 h-4" />
                  <span className="text-sm">Restaurants</span>
                </NavLink>
                <NavLink
                  to="/settings"
                  className={({ isActive }) =>
                    `flex items-center gap-2 px-3 py-1.5 rounded-full transition-all duration-300 whitespace-nowrap ${
                      isActive
                        ? 'bg-gradient-to-r from-blue-500/20 to-purple-500/20 backdrop-blur-md border border-blue-500/30 shadow-lg shadow-blue-500/20'
                        : 'hover:bg-white/20 dark:hover:bg-white/10 hover:backdrop-blur-md hover:shadow-md'
                    }`
                  }
                >
                  <Settings className="w-4 h-4" />
                  <span className="text-sm">Settings</span>
                </NavLink>
              </div>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route path="/" element={<HomePage filters={filters} />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
