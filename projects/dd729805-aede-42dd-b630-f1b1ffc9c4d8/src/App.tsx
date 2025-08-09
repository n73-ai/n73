function App() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-500 via-purple-600 to-pink-500 flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-white mb-4 drop-shadow-lg">
          Hello World!
        </h1>
        <p className="text-xl text-white/90 mb-8">
          Welcome to your React TypeScript app
        </p>
        <div className="bg-white/20 backdrop-blur-sm rounded-lg p-6 border border-white/30">
          <p className="text-white font-medium">
            Built with React, TypeScript & Tailwind CSS
          </p>
        </div>
      </div>
    </div>
  );
}

export default App;
