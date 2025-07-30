type SpinnerProps = {
  color?: string; 
};

export default function Spinner({ color = "#fff" }: SpinnerProps) {

  const barStyle = {
    backgroundColor: color,
  };

  return (
    <div className="spinner">
      <div className="bar1" style={barStyle}></div>
      <div className="bar2" style={barStyle}></div>
      <div className="bar3" style={barStyle}></div>
      <div className="bar4" style={barStyle}></div>
      <div className="bar5" style={barStyle}></div>
      <div className="bar6" style={barStyle}></div>
      <div className="bar7" style={barStyle}></div>
      <div className="bar8" style={barStyle}></div>
    </div>
  );
}
