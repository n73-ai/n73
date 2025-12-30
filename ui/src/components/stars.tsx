import { useState, useEffect } from "react";
import Spinner from "./spinner";
import { AlertCircleIcon } from "lucide-react";

interface Star {
  id: number;
  x: number;
  y: number;
  size: number;
  opacity: number;
  minOpacity: number;
  maxOpacity: number;
  twinkleDelay: number;
  duration: number;
  appearDelay: number;
}

const StarryBackground = ({
  isIframeError,
  status,
}: {
  isIframeError: boolean;
  status: string;
}) => {
  const [stars, setStars] = useState<Star[]>([]);
  const [mounted, setMounted] = useState<boolean>(false);

  useEffect(() => {
    const generateStars = (): void => {
      const newStars: Star[] = [];
      const numberOfStars = 450;

      for (let i = 0; i < numberOfStars; i++) {
        const baseOpacity = Math.random() * 0.4 + 0.6;
        const animationDelay = Math.random() * 8;
        const animationDuration = 6 + Math.random() * 4;
        const appearDelay = i * 0.08 + Math.random() * 4;

        newStars.push({
          id: i,
          x: Math.random() * 100,
          y: Math.random() * 100,
          size: Math.random() * 3 + 1,
          opacity: baseOpacity,
          minOpacity: baseOpacity * 0.7,
          maxOpacity: baseOpacity,
          twinkleDelay: animationDelay,
          duration: animationDuration,
          appearDelay: appearDelay,
        });
      }

      setStars(newStars);
    };

    generateStars();

    const timer = setTimeout(() => {
      setMounted(true);
    }, 500);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="relative w-full h-screen overflow-hidden">
      <style>{`
        @keyframes twinkle {
          0% {
            opacity: calc(
              var(--min-opacity) + (var(--max-opacity) - var(--min-opacity)) *
                0.5 + (var(--max-opacity) - var(--min-opacity)) * 0.5 *
                sin(0deg)
            );
            transform: scale(1);
          }
          25% {
            opacity: var(--max-opacity);
            transform: scale(1.1);
          }
          50% {
            opacity: calc(
              var(--min-opacity) + (var(--max-opacity) - var(--min-opacity)) *
                0.5 + (var(--max-opacity) - var(--min-opacity)) * 0.5 *
                sin(180deg)
            );
            transform: scale(1);
          }
          75% {
            opacity: var(--max-opacity);
            transform: scale(1.1);
          }
          100% {
            opacity: calc(
              var(--min-opacity) + (var(--max-opacity) - var(--min-opacity)) *
                0.5 + (var(--max-opacity) - var(--min-opacity)) * 0.5 *
                sin(360deg)
            );
            transform: scale(1);
          }
        }

        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: scale(0.3);
          }
          to {
            opacity: var(--target-opacity);
            transform: scale(1);
          }
        }

@keyframes skyMovement {
  0% { transform: translateX(100vw) translateY(-3vh) rotate(1deg); }
  10% { transform: translateX(85vw) translateY(2vh) rotate(-0.5deg); }
  20% { transform: translateX(70vw) translateY(-1vh) rotate(0.8deg); }
  30% { transform: translateX(55vw) translateY(3vh) rotate(-0.3deg); }
  40% { transform: translateX(40vw) translateY(-2vh) rotate(0.6deg); }
  50% { transform: translateX(25vw) translateY(2.5vh) rotate(-0.7deg); }
  60% { transform: translateX(10vw) translateY(-1.5vh) rotate(0.4deg); }
  70% { transform: translateX(-5vw) translateY(3vh) rotate(-0.8deg); }
  80% { transform: translateX(-20vw) translateY(-2vh) rotate(0.5deg); }
  90% { transform: translateX(-35vw) translateY(2vh) rotate(-0.6deg); }
  100% { transform: translateX(-100vw) translateY(-3vh) rotate(1deg); }
}

        .sky-container {
          animation: skyMovement 260s linear infinite;
          width: 300vw;
          height: 110%;
          position: absolute;
          top: -5%;
          left: -100vw;
        }

        .star {
          opacity: 0;
          animation: twinkle 4s ease-in-out infinite;
        }

        .star.visible {
          animation:
            fadeIn 2.5s ease-out forwards,
            twinkle 4s ease-in-out infinite;
        }
      `}</style>
      <div className="sky-container">
        {stars.map((star) => (
          <div
            key={star.id}
            className={`absolute rounded-full bg-zinc-400 star ${mounted ? "visible" : ""}`}
            style={
              {
                left: `${star.x}%`,
                top: `${star.y}%`,
                width: `${star.size}px`,
                height: `${star.size}px`,
                "--min-opacity": star.minOpacity,
                "--max-opacity": star.maxOpacity,
                "--target-opacity": star.maxOpacity,
                animationDelay: `${star.appearDelay}s, ${star.twinkleDelay + star.appearDelay + 2.5}s`,
                animationDuration: `2.5s, ${star.duration}s`,
                boxShadow: `0 0 ${star.size * 1.5}px rgba(255, 255, 255, ${star.opacity * 0.3})`,
              } as React.CSSProperties
            }
          />
        ))}

        {stars.slice(0, 20).map((star) => (
          <div
            key={`bright-${star.id}`}
            className={`absolute rounded-full bg-zinc-200 star ${mounted ? "visible" : ""}`}
            style={
              {
                left: `${star.x}%`,
                top: `${star.y}%`,
                width: `${star.size + 1}px`,
                height: `${star.size + 1}px`,
                "--min-opacity": star.minOpacity,
                "--max-opacity": star.maxOpacity,
                "--target-opacity": star.maxOpacity,
                animationDelay: `${star.appearDelay + 1}s, ${star.twinkleDelay + star.appearDelay + 3.5}s`,
                animationDuration: `2.5s, ${star.duration + 1}s`,
                boxShadow: `0 0 ${star.size * 2}px rgba(255, 255, 255, 0.4)`,
              } as React.CSSProperties
            }
          />
        ))}
      </div>

      <div className="relative z-10 flex items-center justify-center h-full">
        <div className="text-center text-white">
          {status == "new_error" && (
            <div className="text-xl flex items-center gap-2 text-destructive py-[30px]">
              <AlertCircleIcon />
              Code compilation failed
            </div>
          )}

          {isIframeError && status != "new_pending" && status != "new_error" && (
            <div className="text-xl flex items-center gap-2 text-destructive py-[30px]">
              <AlertCircleIcon />
              Your website looks it's not online
            </div>
          )}

          {status == "new_pending" && (
            <div className="text-xl flex items-center gap-2 text-muted-foreground py-[30px]">
              <Spinner />
              Building your project
            </div>
          )}

        </div>
      </div>
    </div>
  );
};

export default StarryBackground;
