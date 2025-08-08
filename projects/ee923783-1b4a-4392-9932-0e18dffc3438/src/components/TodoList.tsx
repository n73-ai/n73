import { useTodo } from '../context/TodoContext';
import { TodoItem } from './TodoItem';

export function TodoList() {
  const { todos } = useTodo();

  if (todos.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-zinc-400 text-lg">No tasks yet. Add one to get started!</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {todos.map((todo) => (
        <TodoItem key={todo.id} todo={todo} />
      ))}
    </div>
  );
}