import React from 'react';
import ReactDOM from 'react-dom/client';
import './App.css';
import App from './App.js';
import reportWebVitals from './reportWebVitals';
import { BrowserRouter } from 'react-router-dom';
import { configureStore } from "@reduxjs/toolkit";
import userSlice from './redux/userSlice.js';
import { Provider } from "react-redux";
import dataSlice from './redux/dataSlice';
import { WebSocketProvider } from './contexts/WebSocketContext';
import { library } from '@fortawesome/fontawesome-svg-core';
import { fas } from '@fortawesome/free-solid-svg-icons';

const root = ReactDOM.createRoot(document.getElementById('root'));

library.add(fas);

const store = configureStore({
  reducer: {
    user: userSlice,
    data: dataSlice
  },
});

root.render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <WebSocketProvider>
          <App />
        </WebSocketProvider>
      </BrowserRouter>
    </Provider>
  </React.StrictMode >
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
