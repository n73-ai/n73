import { useState, FormEvent } from 'react';
import { useTodo } from '../context/TodoContext';

export function TodoForm() {
  const [text, setText] = useState('');
  const { addTodo } = useTodo();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (text.trim()) {
      addTodo(text.trim());
      setText('');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-6">
      <div className="flex gap-2">
        <input
          type="text"
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Add a new task..."
          className="flex-1 px-4 py-3 bg-zinc-800 text-zinc-200 rounded-lg border border-zinc-700 focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-500/20 placeholder-zinc-500"
        />
        <button
          type="submit"
          className="px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={!text.trim()}
        >
          Add Task
        </button>
      </div>
    </form>
  );
}