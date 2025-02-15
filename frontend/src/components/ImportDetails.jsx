import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import {Link} from "react-router-dom";
import { faSquarePlus, faBan, faRetweet, faEquals, faTriangleExclamation } from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";


const ImportDetails = () => {
    const { id } = useParams();
    const [importEntry, setImportEntry] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const fetchImport = async () => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/imports/${id}`);
            setImportEntry(response.data);
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
        fetchImport();
    }, [id]);

    if (loading) return <div className="text-center mt-5"><div className="spinner-border" role="status"></div></div>;

    return (
        <div className="container mt-5">
            {error && <div className="alert alert-danger" role="alert">
                {error}
            </div>}
            <div className="card shadow mb-3">
                <div className="card-header">
                    <h2 className="card-title">
                        Import {importEntry.id}
                    </h2>
                </div>
                <div className="card-body">
                    <ul className="list-group list-group-flush">
                        <li className="list-group-item d-flex">
                            <span><strong>Channel</strong></span>
                            <span className="ms-3">
                                <Link to={`/channels/${importEntry.channel_id}`}>{importEntry.channel_name}</Link> ({importEntry.integration})
                            </span>
                        </li>
                        <li className="list-group-item d-flex">
                            <span><strong>Status</strong></span> <span className="ms-3">{importEntry.status}</span>
                        </li>
                        <li className="list-group-item d-flex">
                            <span><strong>Start</strong></span> <span className="ms-3">{importEntry.started_at}</span>
                        </li>
                        <li className="list-group-item d-flex">
                            <span><strong>End</strong></span> <span className="ms-3">{importEntry.ended_at}</span>
                        </li>
                    </ul>
                    <div className="table-responsive mt-3">
                        <table className="table table-striped">
                            <thead>
                            <tr>
                                <th scope="col">New</th>
                                <th scope="col">Updated</th>
                                <th scope="col">No Changes</th>
                                <th scope="col">Missing</th>
                                <th scope="col">Failed</th>
                                <th scope="col">Total</th>
                            </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td><span className="me-1" title="new"><FontAwesomeIcon icon={faSquarePlus} /> {importEntry.new_jobs}</span></td>
                                    <td><span className="me-1" title="new"><FontAwesomeIcon icon={faRetweet} /> {importEntry.updated_jobs}</span></td>
                                    <td><span className="me-1" title="new"><FontAwesomeIcon icon={faEquals} /> {importEntry.no_change_jobs}</span></td>
                                    <td><span className="me-1" title="new"><FontAwesomeIcon icon={faBan} /> {importEntry.missing_jobs}</span></td>
                                    <td><span className="me-1" title="new"><FontAwesomeIcon icon={faTriangleExclamation} /> {importEntry.failed_jobs}</span></td>
                                    <td>{importEntry.total_jobs}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default ImportDetails;