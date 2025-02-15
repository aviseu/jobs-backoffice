import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import { useNavigate } from "react-router-dom";


const ChannelDetails = () => {
    const { id } = useParams();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [updating, setUpdating] = useState(false);
    const [name, setName] = useState("");
    const navigate = useNavigate();

    const fetchChannel = async () => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/channels/${id}`);
            if (response.data.name) {
                setName(response.data.name);
            } else {
                setError('Unexpected non name on response')
                console.error('Unexpected non name on response:', data);
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

    useEffect(() => {
        fetchChannel();
    }, [id]);

    const handleSubmit = async (event) => {
        event.preventDefault();
        setUpdating(true);
        setError(null);
        try {
            const response = await axios.patch(`${import.meta.env.VITE_BACKEND_URL}/api/channels/${id}`, {name: name});
            setTimeout(() => navigate("/channels/" + id), 1);
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
            setUpdating(false);
        }
    }

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
                        <button type="submit" className="btn btn-primary">Update</button>
                    </form>
                </div>
            </div>
        </div>
    )
}

export default ChannelDetails;