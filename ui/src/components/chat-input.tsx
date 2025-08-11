import React, { useState, useRef, useEffect } from 'react';
import { Send, Paperclip, Square } from 'lucide-react';

const ChatInput = () => {
  const [prompt, setPrompt] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const textareaRef = useRef(null);

  // Auto-resize textarea
  useEffect(() => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = 'auto';
      textarea.style.height = Math.min(textarea.scrollHeight, 200) + 'px';
    }
  }, [prompt]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (prompt.trim() && !isGenerating) {
      console.log('Sending message:', prompt);
      setIsGenerating(true);
      // Simulate API call
      setTimeout(() => {
        setIsGenerating(false);
        setPrompt('');
      }, 2000);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleStop = () => {
    setIsGenerating(false);
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <form onSubmit={handleSubmit} className="relative">
        <div className="relative bg-white border border-gray-200 rounded-2xl shadow-sm hover:shadow-md transition-shadow duration-200 focus-within:shadow-md focus-within:border-gray-300">
          {/* Textarea */}
          <textarea
            ref={textareaRef}
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Ask n73 to build something cool..."
            className="w-full px-4 py-3 pr-20 text-gray-900 placeholder-gray-500 bg-transparent border-none rounded-2xl resize-none focus:outline-none focus:ring-0 min-h-[52px] max-h-[200px] leading-6"
            rows={1}
            disabled={isGenerating}
          />
          
          {/* Bottom toolbar */}
          <div className="flex items-center justify-between px-3 pb-3">
            <div className="flex items-center space-x-2">
              {/* Attach button */}
              <button
                type="button"
                className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-150"
                disabled={isGenerating}
              >
                <Paperclip className="w-4 h-4" />
              </button>
            </div>
            
            {/* Send/Stop button */}
            <div className="flex items-center">
              {isGenerating ? (
                <button
                  type="button"
                  onClick={handleStop}
                  className="p-2 text-white bg-gray-600 hover:bg-gray-700 rounded-lg transition-colors duration-150 flex items-center space-x-1"
                >
                  <Square className="w-4 h-4 fill-current" />
                </button>
              ) : (
                <button
                  type="submit"
                  disabled={!prompt.trim()}
                  className={`p-2 rounded-lg transition-all duration-150 flex items-center space-x-1 ${
                    prompt.trim()
                      ? 'text-white bg-black hover:bg-gray-800 shadow-sm hover:shadow-md'
                      : 'text-gray-400 bg-gray-100 cursor-not-allowed'
                  }`}
                >
                  <Send className="w-4 h-4" />
                </button>
              )}
            </div>
          </div>
        </div>
        
        {/* Character counter (optional) */}
        {prompt.length > 0 && (
          <div className="text-xs text-gray-400 mt-2 text-right">
            {prompt.length} characters
          </div>
        )}
      </form>
      
      {/* Demo messages */}
      {isGenerating && (
        <div className="mt-4 p-3 bg-gray-50 rounded-lg border">
          <div className="flex items-center space-x-2">
            <div className="animate-spin w-4 h-4 border-2 border-gray-300 border-t-gray-600 rounded-full"></div>
            <span className="text-sm text-gray-600">n73 is thinking...</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default ChatInput;
