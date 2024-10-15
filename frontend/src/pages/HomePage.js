import SummaryScores from '../components/SummaryScores';
import CalibrationChart from '../components/CalibrationChart';
import './HomePage.css';
import { Link } from 'react-router-dom';



function HomePage() {
return (
    <div className='home-page'>
        <section className='project-info'>
            <h1>Forecasting to understand Reality</h1>
            <p>I forecast to improve my models of the world. On this website, I'll display my current and previous forecasts along with
                my track record.</p>    
        </section>
        <section className='project-info'>
            <h1>Scores</h1>
            <p>Forecasts are scored on their accuracy. The closer each score is to 0, the better. For more information, see 
                <a href='/faq' className='special-link'> FAQ</a>. Click on a datapoint to see information about the forecast.
            </p>
            <SummaryScores></SummaryScores>
        </section>
        <section className='project-info'>
            <h1>Calibration</h1>
            <p>Forecasts are not always right, but I aim to have calibrated forecasts. Meaning that when I forecast
                50% probability, it happens around 50% of the time. Below is my calibration.
            </p>
            <CalibrationChart></CalibrationChart>
        </section>
        <h2>Focus Areas</h2>
        <section className='focus-grid'>
            <div>
                <Link to={'/questions/category/ai'} className='cat-link'>Artifical Intelligence</Link>
                <p>Current and future AI models have the potential to radically change society and humanity for both better and worse. These forecasts
                explore future AI technology and impacts on humanity.
                </p>
                <Link to={'/questions/category/ai'} className='text-link'>See AI Forecasts → </Link>
            </div>
            <div>
                <Link to={'/questions/category/economy'} className='cat-link'>Economy</Link>
                <p>The economy affects peoples lives daily. These forecasts model economic development to help aid economic decision-making.</p>
                <Link to={'/questions/category/economy'} className='text-link'>See Economy Forecasts → </Link>
            </div>
            <div>
                <Link to={'/questions/category/politics'} className='cat-link'>Politics</Link>
                <p>Political developments and changes affect both country and world developments. These forecasts anticipate political developments
                to understand future changes.
                </p>    
                <Link to={'/questions/category/politics'} className='text-link'>See Politics Forecasts → </Link>
            </div>
            <div>
                <Link to={'/questions/category/x-risk'} className='cat-link'>Existential Risks</Link>
                <p>As humanity becomes more and more technologically mature, our technologies are becoming strong enough to kill us all. 
                 <a href='/questions/category/nuclear' className='special-link'>Nuclear weapons</a>, <a href='/questions/category/nuclear' className='special-link'>Artificial Intelligence</a>
                   and Biological weapons are examples of such technologies. 
                </p>
                <Link to={'/questions/category/x-risk'} className='text-link'>See X-Risk Forecasts → </Link>
            </div>
            <div>
                <Link to={'/questions'} className='cat-link'>Other Categories</Link>
                <p>Apart from the focus areas, I also forecast personal questions, Sweden-specific topics, Sports, among many other types of questions.</p>
                <Link to={'/questions'} className='text-link'>See All Forecasts → </Link>
            </div>
        </section>
    </div>
);
};
export default HomePage; 