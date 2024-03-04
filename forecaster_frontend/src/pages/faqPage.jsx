import './HomePage.css';
import { Link } from 'react-router-dom';
import Latex from 'react-latex';
import 'katex/dist/katex.min.css'; // Import KaTeX CSS

function FaqPage() {
    const brierScoreString='BS=\\frac{1}{N} \\sum^N_{t=1} (f_t-o_t)^2'
    const logScoreString='L({\\bf{r}}, i)=\\log_b(r_i)'
    const logScoreString2='L({\\bf{r}}, i)=x\\log_b(p) + (1-x)\\log_b(1-p)'
    return (
        <div className='home-page'>
            <section className='project-info'>
                <h3>FAQ 1: Why are you doing this?</h3>
                <p>I want to make my beliefs match reality as much as possible. To make this happen, this public forecasting project is 
                    a vital part. Now, my track record is visible to anyone—giving me an incentive to perform at my best in my forecasts.
                    </p>   
                <p>Eventually, I would like to add comments and the ability for readers to challenge my beliefs. Until then, you can 
                    reach out to me with a question or challenge.</p> 
            </section>
            <section className='project-info'>
                <h3>FAQ 2: How are the scores you use defined?</h3>
                <p>I use three different scoring rules: (i) Brier score, (ii) Binary Log score, and (iii) Natural Log score.
                    They are defined as follows. </p>
                <p>(i) Brier score, when applied to one-dimensional forecasts, is the mean squared error of a forecast. This score is 
                    minimized, so a score closer to 0 is better. A maximally uncertain prediction of 50% would yield a Brier score of 0.25.                
                    </p>
                    <Latex>{`$${brierScoreString}$`}</Latex> 
                <p>(ii) Binary Log Score and (iii) Natural Log score are equivalent to the negative of Surprisal, a 
                    commonly used scoring criterion in Bayesian inference. Compared to Brier score, these scores are harsher on poor
                    predictions that are closer to 0 and 1. These scores are maximized, so a prediction scoring -0.10 is better than a 
                    prediction scoring -0.8.
                </p>
                <Latex>{`$${logScoreString}$`}</Latex> 
                <p>This is a local strictly proper scoring rule under any base larger than 1. Given a binary prediction, the logarithmic 
                    scoring rules can be rewritten as:
                </p>
                <Latex>{`$${logScoreString2}$`}</Latex> 
                <p>Where x is the outcome of the event and p is the expressed probability. </p>
            </section>
            <section className='project-info'>
                <h3>FAQ 3: What is a proper scoring rule?</h3>
                <p>A proper scoring rule is rule which, when implemented, incentivizes forecasters to report their True beliefs when forecasting.
                    Mathematically, this means that a forecaster which maximizes reward will report the true probability distribution when they know the true
                    probability distribution.</p>    
            </section>
            <section className='project-info'>
                <h3>FAQ 4: How can I contact you?</h3>
                <p>You can reach me on email: bayesmaxxing@gmail.com</p>    
            </section>
            <section className='project-info'>
                <h3>FAQ 5: What are your plans for this project?</h3>
                <p>As a next step, I would like to be able to present statistical models that I've used to make my forecasts.
                    Mostly, I imagine this being a link to a Google Colab notebook, but eventually this could turn into a more advanced feature.</p>    
                <p>Then, I would like to add the possibility for readers to comment and suggest questions for me to forecast.</p>
            </section>
            <section className='project-info'>
                <h3>FAQ 6: Can I also forecast?</h3>
                <p>Of course, anyone can do forecasting.</p>    
            </section>
            <section className='project-info'>
                <h3>FAQ 7: Why are many of the questions and predictions made on 2023-10-16?</h3>
                <p>Before starting the project, I did not track the timing of my forecasts. To have a time of creation for all questions and predictions, 
                    I set the time to the date I migrated my forecasts to the new database. </p>
            </section>
        </div>
    );
    };
    export default FaqPage; 