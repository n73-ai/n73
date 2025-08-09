import { useState } from 'react';

function App() {
  const [isDarkMode, setIsDarkMode] = useState(true);

  const toggleTheme = () => {
    setIsDarkMode(!isDarkMode);
  };

  return (
    <div className={`min-h-screen flex items-center justify-center transition-all duration-500 ${
      isDarkMode 
        ? 'bg-gradient-to-br from-gray-900 via-gray-800 to-black' 
        : 'bg-gradient-to-br from-yellow-200 via-pink-200 to-blue-200'
    }`}>
      <div className="text-center relative">
        {/* Theme Toggle Button */}
        <button
          onClick={toggleTheme}
          className={`absolute -top-16 right-0 p-3 rounded-full backdrop-blur-sm border transition-all duration-300 hover:scale-110 ${
            isDarkMode 
              ? 'bg-gray-800/60 border-gray-700/50 text-white hover:bg-gray-700/70' 
              : 'bg-white/20 border-white/30 text-white hover:bg-white/30'
          }`}
          aria-label="Toggle theme"
        >
          {isDarkMode ? '☀️' : '🌙'}
        </button>

        <h1 className={`text-6xl font-bold mb-4 drop-shadow-lg transition-colors duration-300 ${
          isDarkMode ? 'text-white' : 'text-gray-800'
        }`}>
          Hello World! 👋
        </h1>
        
        <p className={`text-xl max-w-md mx-auto leading-relaxed transition-colors duration-300 ${
          isDarkMode ? 'text-gray-300' : 'text-gray-700'
        }`}>
          Welcome to your React + TypeScript + Tailwind CSS application
        </p>
        
        <div className="mt-8">
          <div className={`inline-block backdrop-blur-sm rounded-full px-6 py-3 border transition-all duration-300 shadow-lg ${
            isDarkMode 
              ? 'bg-gray-800/60 border-gray-700/50' 
              : 'bg-white/40 border-gray-300/50'
          }`}>
            <span className={`font-medium transition-colors duration-300 ${
              isDarkMode ? 'text-white' : 'text-gray-800'
            }`}>
              ✨ Built with React & Tailwind ✨
            </span>
          </div>
        </div>
        
        <div className="mt-6">
          <p className={`text-sm transition-colors duration-300 ${
            isDarkMode ? 'text-white/70' : 'text-gray-600'
          }`}>
            {isDarkMode ? 'Dark Theme Active' : 'Light Theme Active'}
          </p>
        </div>
      </div>
    </div>
  );
}

export default App;
