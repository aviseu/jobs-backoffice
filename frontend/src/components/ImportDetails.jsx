import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import {Link, useLocation } from "react-router-dom";
import { faSquarePlus, faPlus, faBan, faRetweet, faEquals, faQuestion, faCircleQuestion, faFolderPlus } from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";


const ImportDetails = () => {
    const { id } = useParams();
    const [importEntry, setImportEntry] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const location = useLocation();

    const fetchImport = async () => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/imports/${id}`);
            setImportEntry(response.data);
            if (response.data.status !== "completed" && response.data.status !== "failed" && window.location.pathname === "/imports/" + id && new Date() - Date.parse(response.data.started_at) < 5 * 60 * 1000) {
                setTimeout(fetchImport, 500);
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
        fetchImport();
    }, [id, location.pathname]);

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
                    {importEntry.status === "failed" && (
                        <div className="alert alert-danger" role="alert">
                            <h4 className="alert-heading">Import Failed</h4>
                            <hr/>
                            <p>{importEntry.error || "An error occurred during the import."}</p>
                        </div>
                    )}
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
                                <th scope="col">Total</th>
                                <th scope="col">Errors</th>
                                <th scope="col">Missing</th>
                                <th scope="col">P. Missing</th>
                                <th scope="col">P. Info</th>
                                <th scope="col">P. Late</th>
                            </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td><span className="me-1" title="new"><FontAwesomeIcon icon={faPlus} /> {importEntry.new_jobs}</span></td>
                                    <td><span className="me-1" title="updated"><FontAwesomeIcon icon={faRetweet} /> {importEntry.updated_jobs}</span></td>
                                    <td><span className="me-1" title="not changed"><FontAwesomeIcon icon={faEquals} /> {importEntry.no_change_jobs}</span></td>
                                    <td>{importEntry.total_jobs}</td>
                                    <td><span className="me-1" title="failed"><FontAwesomeIcon icon={faBan} /> {importEntry.errors}</span></td>
                                    <td><span className="me-1" title="missing"><FontAwesomeIcon icon={faQuestion} /> {importEntry.missing_jobs}</span></td>
                                    <td><span className="me-1" title="missing published"><FontAwesomeIcon icon={faCircleQuestion} /> {importEntry.missing_published}</span></td>
                                    <td><span className="me-1" title="published"><FontAwesomeIcon icon={faSquarePlus} /> {importEntry.published}</span></td>
                                    <td><span className="me-1" title="late published"><FontAwesomeIcon icon={faFolderPlus} /> {importEntry.late_published}</span></td>
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