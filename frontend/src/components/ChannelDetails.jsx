import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axios from "axios";
import {Link} from "react-router-dom";


const ChannelDetails = () => {
    const { id } = useParams();
    const [channel, setChannel] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();
    const [updating, setUpdating] = useState(false);

    const fetchChannel = async () => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/channels/${id}`);
            setChannel(response.data);
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

    const scheduleImport = async (event) => {
        event.preventDefault();
        setUpdating(true);
        setError(null);
        try {
            const response = await axios.put(`${import.meta.env.VITE_BACKEND_URL}/api/channels/${id}/schedule`);
            setTimeout(() => navigate("/imports/" + response.data.id), 0);
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

    const changeStatus = async (action, event) => {
        event.preventDefault();
        setUpdating(true);
        setError(null);
        try {
            const response = await axios.put(`${import.meta.env.VITE_BACKEND_URL}/api/channels/${id}/${action}`);
            await fetchChannel();
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
        <div className="container mt-5">
            <div className={`card border-2 shadow ${ channel.status === "active" ? "border-success" : "" }`}>
                <div className="card-body">
                    <h2 className="card-title mb-5">
                        {channel.name}

                        {channel.status === "active" && (
                            <button className="btn btn-sm btn-danger float-end"
                                    onClick={(event) => changeStatus("deactivate", event)} disabled={updating}>
                                {updating ?
                                    <span className="spinner-border spinner-border-sm"></span> : "Deactivate"}
                            </button>
                        )}

                        {channel.status === "inactive" && (
                            <button className="btn btn-sm btn-success float-end "
                                    onClick={(event) => changeStatus("activate", event)} disabled={updating}>
                                {updating ? <span className="spinner-border spinner-border-sm"></span> : "Activate"}
                            </button>
                        )}

                        <button className="btn btn-sm btn-warning float-end me-2"
                                onClick={(event) => scheduleImport(event)} disabled={updating}>
                            {updating ?
                                <span className="spinner-border spinner-border-sm"></span> : "Import"}
                        </button>

                        <Link className="btn btn-sm btn-primary float-end me-2" role="button" to={"/channels/"+id+"/update"}>Update</Link>
                    </h2>
                    <h6 className="mb-3">Integration: {channel.integration}</h6>
                </div>
            </div>
        </div>
    )
}

export default ChannelDetails;