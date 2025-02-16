import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from "react-router-dom";

const ChannelCreate = () => {
    const [name, setName] = useState('');
    const [integration, setIntegration] = useState('');
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    const [options, setOptions] = useState([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const fetchOptions = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/integrations`);
                const data = response.data.integrations;
                if (Array.isArray(data)) {
                    setOptions(data);
                } else {
                    setError('Unexpected non array response')
                    console.error('Expected an array but got:', data);
                }
            } catch (err) {
                console.log(err)
                if (err.response) {
                    setError(err.response.data.error.message || "Submission failed. Please check your input.");
                } else if (err.request) {
                    setError("Network error. Please check your connection.");
                } else {
                    setError("An unexpected error occurred.");
                }
            } finally {
                setLoading(false);
            }
        };

        fetchOptions();
    }, []);

    const handleChange = (event) => {
        setIntegration(event.target.value);
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        setLoading(true);
        setError(null);
        try {
            const response = await axios.post(`${import.meta.env.VITE_BACKEND_URL}/api/channels`, {name, integration})
            setTimeout(() => navigate("/channels/" + response.data.id), 0);
        } catch (err) {
            console.log(err)
            if (err.response) {
                setError(err.response.data.error.message || "Submission failed. Please check your input.");
            } else if (err.request) {
                setError("Network error. Please check your connection.");
            } else {
                setError("An unexpected error occurred.");
            }
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="text-center mt-5"><div className="spinner-border" role="status"></div></div>;

    return (
        <div>
            <div className="row justify-content-md-center mt-5">
                <div className="col-6">
                    <h1 className="h2">Create Channel</h1>
                    {error && <div className="alert alert-danger" role="alert">
                        {error}
                    </div>}
                    <form onSubmit={handleSubmit}>
                        <div className="mb-3">
                            <label htmlFor="exampleFormControlInput1" className="form-label">Name</label>
                            <input type="text" className="form-control" value={name}
                                   onChange={(e) => setName(e.target.value)}/>
                        </div>
                        <div className="mb-3">
                            <label htmlFor="exampleFormControlTextarea1" className="form-label">Integration</label>
                            <select value={integration} onChange={handleChange} className="form-select">
                                <option value="">&nbsp;</option>
                                {options.map((option) => (
                                    <option key={option} value={option}>
                                        {option}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <button type="submit" className="btn btn-primary">Create</button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ChannelCreate;
