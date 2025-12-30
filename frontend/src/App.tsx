import { useState, lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom';
import { Home, Settings, Loader2 } from 'lucide-react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useTheme } from './hooks/useTheme';
import { ThemeToggle } from './components/ThemeToggle';
import { GlobalSearch } from './components/GlobalSearch';
import { useCategories, useFoodTypes } from './hooks/useApi';
import { RestaurantFilters } from './services/api';
import { ErrorBoundary } from './components/ErrorBoundary';
import { ToastProvider } from './hooks/useToast';

// Lazy load page components for code splitting
const HomePage = lazy(() => import('./pages/HomePage').then(m => ({ default: m.HomePage })));
const SettingsPage = lazy(() => import('./pages/SettingsPage').then(m => ({ default: m.SettingsPage })));

// Loading fallback component
const PageLoader = () => (
  <div className="flex items-center justify-center h-64">
    <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
  </div>
);

// Create QueryClient instance with optimized defaults and error handling
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false, // Don't refetch on window focus
      retry: 1, // Retry failed requests once
      staleTime: 5 * 60 * 1000, // 5 minutes default
    },
    mutations: {
      retry: 0, // Don't retry mutations by default
    },
  },
});

function AppContent() {
  const { isDark, toggleTheme } = useTheme();
  const [filters, setFilters] = useState<RestaurantFilters>({});

  // Use React Query hooks instead of manual fetching
  const { data: categories = [] } = useCategories();
  const { data: foodTypes = [] } = useFoodTypes();

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
          <Suspense fallback={<PageLoader />}>
            <Routes>
              <Route path="/" element={<HomePage filters={filters} />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Routes>
          </Suspense>
        </main>
      </div>
    </BrowserRouter>
  );
}

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <ToastProvider>
          <AppContent />
          <ReactQueryDevtools initialIsOpen={false} />
        </ToastProvider>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}

export default App;
