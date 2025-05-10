const THREE = require("three");
const { GLTFLoader } = require("three/examples/jsm/loaders/GLTFLoader.js");
const { DRACOLoader } = require("three/examples/jsm/loaders/DRACOLoader.js");
const fs = require("fs");
const path = require("path");

// Mock DOM elements that GLTFLoader might expect
global.document = {
    createElementNS: () => ({}),
    createElement: () => ({})
};
global.Image = function() {};
global.XMLHttpRequest = require("xmlhttprequest").XMLHttpRequest;

const modelPath = path.resolve(__dirname, "../models/LittlestTokyo.glb");
const dracoDir = path.resolve(__dirname, "../js/draco/");

const loader = new GLTFLoader();
const dracoLoader = new DRACOLoader();
dracoLoader.setDecoderPath(`file://${dracoDir}/`); // Important: use file:// protocol for local paths
loader.setDRACOLoader(dracoLoader);

fs.readFile(modelPath, (err, data) => {
    if (err) {
        console.error("Error reading GLB file:", err);
        return;
    }
    loader.parse(data.buffer, "", (gltf) => {
        console.log("Model loaded successfully. Mesh names:");
        gltf.scene.traverse((child) => {
            if (child.isMesh) {
                console.log(`- Name: '${child.name}', UUID: ${child.uuid}`);
            }
        });
        // Clean up DracoWorker
        dracoLoader.dispose();
    }, (error) => {
        console.error("Error parsing GLTF model:", error);
        // Clean up DracoWorker
        dracoLoader.dispose();
    });
});

