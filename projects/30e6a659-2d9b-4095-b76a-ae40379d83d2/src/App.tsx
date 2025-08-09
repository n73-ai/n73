function App() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-500 via-purple-600 to-pink-500 flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-white mb-4 drop-shadow-lg">
          Hello World! 👋
        </h1>
        <p className="text-xl text-white/90 max-w-md mx-auto leading-relaxed">
          Welcome to your React + TypeScript + Tailwind CSS application
        </p>
        <div className="mt-8">
          <div className="inline-block bg-white/20 backdrop-blur-sm rounded-full px-6 py-3 border border-white/30">
            <span className="text-white font-medium">✨ Built with React & Tailwind ✨</span>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
