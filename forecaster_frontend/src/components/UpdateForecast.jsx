import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const UpdateForecast = () => {
    let { id } = useParams();
    const [updateData, setUpdateData] = useState({
        point_forecast: '',
        upper_ci: '',
        lower_ci: '', 
        reason: '',
    });
    const [submitStatus,setSubmitStatus] = useState('');

    const handleChange = (e) => {
        const { name, value } = e.target;
        if (['point_forecast', 'upper_ci', 'lower_ci'].includes(name)) {
            // Allow numbers, one decimal point, and handle leading decimal
            const regex = /^-?\d*\.?\d*$/;
            if (value === '' || regex.test(value)) {
                setUpdateData(prevState => ({
                    ...prevState,
                    [name]: value
                }));
            }
        } else {
            setUpdateData(prevState => ({
                ...prevState,
                [name]: value
            }));
        }
    };

    const getDate = () => {
        const currentDate = new Date();
        return currentDate.toISOString().split('T')[0] + ' 00:00:00'
    }

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitStatus('Submitting update...');

        const dataToSubmit = {
            ...updateData,
            point_forecast: parseFloat(updateData.point_forecast),
            upper_ci: parseFloat(updateData.upper_ci),
            lower_ci: parseFloat(updateData.lower_ci),  
            reason: updateData.reason,          
            date_added: getDate(),
            forecast: parseInt(id)
        };

        try {
            const response = await fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecast_points/?forecast=${id}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
                }, 
                body: JSON.stringify(dataToSubmit)
            });
            
            if (response.ok) {
                setSubmitStatus('Update added');
                setUpdateData({point_forecast: '', upper_ci:'',
                    lower_ci: '',reason: ''
                });
            } else {
                setSubmitStatus('Update not added. Try again.');    
            } 
        } catch (error) {
            console.error('Error:',error)
            setSubmitStatus('An error occurred. Try again.')
        }
    };
     
    return (
        <div>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label>Point forecast</label>
              <input
                type="text"
                id="point_forecast"
                name="point_forecast"
                value={updateData.point_forecast}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label>Upper confidence interval</label>
              <input
                type="text"
                id="upper_ci"
                name="upper_ci"
                value={updateData.upper_ci}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label>Lower confidence interval</label>
              <input
                type="text"
                id="lower_ci"
                name="lower_ci"
                value={updateData.lower_ci}
                onChange={handleChange}
                required
              ></input>
            </div>
            <div>
              <label >Reason for update</label>
              <textarea
                id="reason"
                name="reason"
                value={updateData.reason}
                onChange={handleChange}
                required
                rows="3"
              ></textarea>
            </div>
            <button
              type="submit"
            >
              Submit
            </button>
          </form>
          {submitStatus && <p >{submitStatus}</p>}
        </div>
      );
    };
    export default UpdateForecast;