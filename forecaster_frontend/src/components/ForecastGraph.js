import './ForecastGraph.css';
import { Line } from 'react-chartjs-2';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
  } from 'chart.js';
  
ChartJS.register(
    CategoryScale, 
    LinearScale, 
    PointElement, 
    LineElement, 
    Title, 
    Tooltip, 
    Legend
);

function ForecastGraph({ data, options }) {
    return (
        <div className="chart-container">
            <Line data={data} options={options} />
        </div>
    );
};
export default ForecastGraph;