import React, {useState, useEffect} from 'react';
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

function CalibrationChart() {
    const [resolutions, setResolutions] = useState([]);
    const [calibrationData, setCalibrationData] = useState([]);
    const [scores, setScores] = useState([]);

    useEffect(() => {
        Promise.all([
            fetch('https://forecasting-389105.ey.r.appspot.com/forecasts?type=resolved', {
                headers : {
                    "Accept" : "application/json"
                }
            }), 
            fetch('https://forecasting-389105.ey.r.appspot.com/forecast-points/latest', {
                headers : {
                    "Accept" : "application/json"
                }
            })
        ])
        .then(async ([resolved, latest]) => {
            const resolvedData = await resolved.json();
            const latestData = await latest.json();
            return [resolvedData, latestData];
        })
        .then(([resolvedData, latestData]) => {
            const combined = resolvedData.map(forecast => {
                const matchingPoint = latestData.find(point => point.forecast_id === forecast.id);
                return {...forecast, latestPoint: matchingPoint || null};
        });
            setResolutions(combined);

            const calibrationData = createCalibrationData(combined);
            setCalibrationData(calibrationData);
        })
        .catch(error => console.error('Error fetching data: ',error));
    }, []);
    
    const createCalibrationData = (data) => {
        const bins = Array(10).fill().map((_, i) => ({
            min: i / 10,
            max: (i + 1) / 10,
            predictions: 0,
            occurrences: 0
        }));

        data.forEach(item => {
            if (item.latestPoint && item.latestPoint.point_forecast !== null) {
                const prediction = item.latestPoint.point_forecast;
                const binIndex = Math.min(Math.floor(prediction * 10), 9);
                bins[binIndex].predictions++;
                if (item.resolution === "1") { 
                    bins[binIndex].occurrences++;
                }
            }
        });

        // Calculate the actual probability for each bin
        return bins.map(bin => ({
            ...bin,
            actualProbability: bin.predictions > 0 ? bin.occurrences / bin.predictions : 0
        }));
    };

    const chartData = {
        labels: calibrationData ? calibrationData.map(bin => `${bin.min * 100}-${bin.max * 100}%`) : [],
        datasets: [
            {
                label: 'Perfect Calibration',
                data: [0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1],
                borderColor: 'rgba(75, 192, 192, 1)',
                borderDash: [5, 5],
                fill: false
            },
            {
                label: 'My Calibration',
                data: calibrationData ? calibrationData.map(bin => bin.actualProbability) : [],
                borderColor: 'rgba(255, 99, 132, 1)',
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                fill: false
            }
        ]
    };

    const chartOptions = {
        scales: {
            x: {
                title: {
                    display: true,
                    text: 'Predicted Probability'
                }
            },
            y: {
                title: {
                    display: true,
                    text: 'Actual Probability'
                },
                min: 0,
                max: 1
            }
        },
        plugins: {
            title: {
                display: true,
                text: 'Calibration Chart'
            }
        }
    };

    return (
        <div>
            {calibrationData && (
                <div>
                    <Line data={chartData} options={chartOptions} />
                </div>
            )}
        </div>
    );
};


export default CalibrationChart;