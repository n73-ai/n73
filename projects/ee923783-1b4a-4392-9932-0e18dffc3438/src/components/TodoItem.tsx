import type { Todo } from '../types/todo';
import { useTodo } from '../context/TodoContext';

interface TodoItemProps {
  todo: Todo;
}

export function TodoItem({ todo }: TodoItemProps) {
  const { toggleTodo, deleteTodo } = useTodo();

  return (
    <div className={`flex items-center justify-between p-4 mb-2 rounded-lg ${
      todo.completed ? 'bg-zinc-800 opacity-75' : 'bg-zinc-700'
    }`}>
      <div className="flex items-center">
        <input
          type="checkbox"
          checked={todo.completed}
          onChange={() => toggleTodo(todo.id)}
          className="w-5 h-5 mr-3 rounded-md text-purple-500 bg-zinc-600 border-zinc-500 focus:ring-purple-500"
        />
        <span className={`text-zinc-200 ${
          todo.completed ? 'line-through text-zinc-400' : ''
        }`}>
          {todo.text}
        </span>
      </div>
      <button
        onClick={() => deleteTodo(todo.id)}
        className="p-1 text-zinc-400 hover:text-red-500 transition-colors"
        aria-label="Delete todo"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="w-5 h-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          />
        </svg>
      </button>
    </div>
  );
}