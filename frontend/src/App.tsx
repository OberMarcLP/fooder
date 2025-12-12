import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom';
import { Home, Tag, Utensils, Lightbulb } from 'lucide-react';
import { useTheme } from './hooks/useTheme';
import { ThemeToggle } from './components/ThemeToggle';
import { HomePage } from './pages/HomePage';
import { CategoriesPage } from './pages/CategoriesPage';
import { FoodTypesPage } from './pages/FoodTypesPage';
import { SuggestionsPage } from './pages/SuggestionsPage';

function App() {
  const { isDark, toggleTheme } = useTheme();

  return (
    <BrowserRouter>
      <div className="min-h-screen">
        <nav className="bg-white dark:bg-gray-800 shadow-lg sticky top-0 z-40">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex items-center justify-between h-16">
              <div className="flex items-center gap-8">
                <span className="text-xl font-bold text-blue-600 dark:text-blue-400">
                  Fooder
                </span>
                <div className="flex gap-1">
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
                    <span className="hidden sm:inline">Restaurants</span>
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
                    <span className="hidden sm:inline">Categories</span>
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
                    <span className="hidden sm:inline">Food Types</span>
                  </NavLink>
                  <NavLink
                    to="/suggestions"
                    className={({ isActive }) =>
                      `flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
                        isActive
                          ? 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400'
                          : 'hover:bg-gray-100 dark:hover:bg-gray-700'
                      }`
                    }
                  >
                    <Lightbulb className="w-4 h-4" />
                    <span className="hidden sm:inline">Suggestions</span>
                  </NavLink>
                </div>
              </div>
              <ThemeToggle isDark={isDark} onToggle={toggleTheme} />
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/categories" element={<CategoriesPage />} />
            <Route path="/food-types" element={<FoodTypesPage />} />
            <Route path="/suggestions" element={<SuggestionsPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
