import { TrackJobApp } from "../wailsjs/go/main/App.js";

console.log("index.js loaded");

window.trackJob = async function() {
    console.log("running trackJob");
    const jobAppStr = document.getElementById("jobApp").value.trim();

    TrackJobApp(jobAppStr);
};