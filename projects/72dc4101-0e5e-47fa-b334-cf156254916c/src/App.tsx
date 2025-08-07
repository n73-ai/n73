function App() {
  return (
    <div className="flex flex-col justify-center items-center bg-white min-h-screen space-y-6">
      <h1 className="text-6xl font-bold text-blue-600">Hello, World!</h1>
      <button 
        className="px-6 py-3 bg-red-600 text-white font-bold rounded-lg hover:bg-red-700 transition-colors"
        onClick={() => alert('Button clicked!')}
      >
        Click Me
      </button>
    </div>
  );
}

export default App;
