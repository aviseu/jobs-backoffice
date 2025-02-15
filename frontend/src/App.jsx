import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import './App.css'
import TopBar from './components/TopBar';
import ChannelList from './components/ChannelList';
import ChannelCreate from './components/ChannelCreate.jsx';
import ChannelUpdate from './components/ChannelUpdate.jsx';
import ChannelDetails from './components/ChannelDetails';
import ImportList from './components/ImportList';
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min";

function App() {
  return (
      <>
          <Router>
              <TopBar/>
              <div className="container">
                  <Routes>
                      <Route path="/" element={<Navigate to="/channels" />} />

                      <Route exact path="/channels" element={<ChannelList/>} />
                      <Route path="/channels/create" element={<ChannelCreate/>} />
                      <Route path="/channels/:id/update" element={<ChannelUpdate/>} />
                      <Route path="/channels/:id" element={<ChannelDetails/>} />

                      <Route path="/imports" element={<ImportList/>} />
                  </Routes>
              </div>
          </Router>
      </>
  );
};

export default App;
