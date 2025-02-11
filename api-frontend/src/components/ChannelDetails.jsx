import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import {Link} from "react-router-dom";


const ChannelDetails = () => {
    const { id } = useParams();
    const [channel, setChannel] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [updating, setUpdating] = useState(false);

    const fetchChannel = async () => {
        try {
            const response = await axios.get(`http://localhost:8080/api/channels/${id}`);
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

    const changeStatus = async (action, event) => {
        event.preventDefault();
        setUpdating(true);
        setError(null);
        try {
            const response = await axios.put(`http://localhost:8080/api/channels/${id}/${action}`);
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
            <div className={`card shadow ${ channel.status === "active" ? "border-success" : "" }`}>
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
                            <button className="btn btn-sm btn-success float-end"
                                    onClick={(event) => changeStatus("activate", event)} disabled={updating}>
                                {updating ? <span className="spinner-border spinner-border-sm"></span> : "Activate"}
                            </button>
                        )}

                        <Link className="btn btn-sm btn-primary float-end me-2" role="button" to={"/"+id+"/update"}>Update</Link>
                    </h2>
                    <h6 className="mb-3">Integration: {channel.integration}</h6>
                </div>
            </div>
        </div>
    )
}

export default ChannelDetails;