// Copyright © 2026 Michael Fero. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

'use strict';

const MAX_BYTES = 10 * 1024 * 1024;
const ACCEPTED_TYPES = new Set(['image/png', 'image/jpeg', 'image/gif']);
const CROP_SIZE = 512;

// State variables; reset on each upload

let currentFile = null;
let naturalWidth = 0;
let naturalHeight = 0;
let cropOffset = 0.5; // 0.0–1.0, default center
let panRange = 0;     // max pan distance in scaled pixels
let isLandscape = false;
let isDragging = false;
let dragStartMouse = 0;
let dragStartPan = 0;
let originalObjectURL = null;
let resultObjectURL = null;

// Element references

const uploadZone  = document.getElementById('upload-zone');
const fileInput   = document.getElementById('file-input');
const errorPanel  = document.getElementById('error-panel');
const errorMsg    = document.getElementById('error-message');
const cropModal   = document.getElementById('crop-modal');
const cropHint    = document.getElementById('crop-hint');
const cropViewport = document.getElementById('crop-viewport');
const cropImage   = document.getElementById('crop-image');
const convertBtn  = document.getElementById('convert-btn');
const statusMsg   = document.getElementById('status-msg');
const resultPanel = document.getElementById('result-panel');
const originalImg = document.getElementById('original-img');
const resultImg   = document.getElementById('result-img');
const downloadBtn = document.getElementById('download-btn');
const resetBtn    = document.getElementById('reset-btn');
const arrowL      = document.getElementById('crop-arrow-l');
const arrowR      = document.getElementById('crop-arrow-r');
const arrowU      = document.getElementById('crop-arrow-u');
const arrowD      = document.getElementById('crop-arrow-d');

// Upload interactions

uploadZone.addEventListener('click', () => fileInput.click());
uploadZone.addEventListener('keydown', (e) => {
  if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); fileInput.click(); }
});

uploadZone.addEventListener('dragover', (e) => {
  e.preventDefault();
  uploadZone.classList.add('dragging');
});
uploadZone.addEventListener('dragleave', () => uploadZone.classList.remove('dragging'));
uploadZone.addEventListener('drop', (e) => {
  e.preventDefault();
  uploadZone.classList.remove('dragging');
  const file = e.dataTransfer.files[0];
  if (file) handleFile(file);
});

fileInput.addEventListener('change', () => {
  if (fileInput.files[0]) handleFile(fileInput.files[0]);
  fileInput.value = '';
});

function handleFile(file) {
  hideError();

  if (!ACCEPTED_TYPES.has(file.type)) {
    showError('FORMAT REJECTED. This binary expected a PNG, JPEG, or GIF. This binary received something else. This binary is not angry. This binary is disappointed, which is different.');
    return;
  }
  if (file.size > MAX_BYTES) {
    showError('FILE REJECTED. This binary does not accept files larger than 10 MB. This binary has standards. This binary would like to keep them. Please reduce your file size and try again.');
    return;
  }

  currentFile = file;
  cropOffset = 0.5;

  // A hidden Image element is the only way to read naturalWidth/naturalHeight from a
  // blob URL before displaying it; dimensions aren't available until onload fires.
  const url = URL.createObjectURL(file);
  const probe = new Image();
  probe.onload = () => {
    naturalWidth = probe.naturalWidth;
    naturalHeight = probe.naturalHeight;
    URL.revokeObjectURL(url);
    showCropModal();
  };
  probe.onerror = () => {
    URL.revokeObjectURL(url);
    showError('READ ERROR. This binary could not process that file as an image. This binary tried. This binary would like credit for trying.');
  };
  probe.src = url;
}

// Crop modal

function showCropModal() {
  const isSquare = naturalWidth === naturalHeight;
  isLandscape = naturalWidth > naturalHeight;

  let scaledW, scaledH;
  if (isSquare) {
    scaledW = CROP_SIZE;
    scaledH = CROP_SIZE;
    panRange = 0;
  } else if (isLandscape) {
    scaledH = CROP_SIZE;
    scaledW = Math.round(naturalWidth * CROP_SIZE / naturalHeight);
    panRange = scaledW - CROP_SIZE;
  } else {
    scaledW = CROP_SIZE;
    scaledH = Math.round(naturalHeight * CROP_SIZE / naturalWidth);
    panRange = scaledH - CROP_SIZE;
  }

  cropImage.style.width = scaledW + 'px';
  cropImage.style.height = scaledH + 'px';
  updateImagePosition();

  const objectURL = URL.createObjectURL(currentFile);
  cropImage.onload = () => URL.revokeObjectURL(objectURL);
  cropImage.src = objectURL;

  cropViewport.classList.toggle('is-square', isSquare);
  cropViewport.classList.toggle('pan-h', !isSquare && isLandscape);
  cropViewport.classList.toggle('pan-v', !isSquare && !isLandscape);

  if (isSquare) {
    cropHint.textContent = 'CONFIRMATION. Your image is already square. This binary acknowledges this is the correct shape. You may proceed.';
  } else if (isLandscape) {
    cropHint.textContent = 'ATTENTION. Your image is not square. Drag left or right to select the area you wish to deactivate. This binary will remember your choice.';
  } else {
    cropHint.textContent = 'ATTENTION. Your image is not square. Drag up or down to select the area you wish to deactivate. This binary will remember your choice.';
  }

  updateArrows();
  convertBtn.disabled = false;
  statusMsg.classList.add('hidden');
  statusMsg.textContent = '';
  uploadZone.classList.add('hidden');
  resultPanel.classList.add('hidden');
  cropModal.classList.remove('hidden');
}

function updateArrows() {
  if (panRange === 0) {
    arrowL.classList.remove('visible');
    arrowR.classList.remove('visible');
    arrowU.classList.remove('visible');
    arrowD.classList.remove('visible');
    return;
  }
  if (isLandscape) {
    arrowL.classList.toggle('visible', cropOffset > 0);
    arrowR.classList.toggle('visible', cropOffset < 1);
    arrowU.classList.remove('visible');
    arrowD.classList.remove('visible');
  } else {
    arrowU.classList.toggle('visible', cropOffset > 0);
    arrowD.classList.toggle('visible', cropOffset < 1);
    arrowL.classList.remove('visible');
    arrowR.classList.remove('visible');
  }
}

function updateImagePosition() {
  const pan = Math.round(cropOffset * panRange);
  cropImage.style.left = (isLandscape ? -pan : 0) + 'px';
  cropImage.style.top  = (isLandscape ? 0 : -pan) + 'px';
}

// Dismiss modal on backdrop click or ESC

cropModal.addEventListener('click', (e) => {
  if (e.target === cropModal) reset();
});

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && !cropModal.classList.contains('hidden')) reset();
});

// Crop viewport panning

cropViewport.addEventListener('mousedown', startPan);
cropViewport.addEventListener('touchstart', (e) => { e.preventDefault(); startPan(e.touches[0]); }, { passive: false });

document.addEventListener('mousemove', onPanMove);
// passive: false is required so e.preventDefault() can suppress native scroll while panning.
document.addEventListener('touchmove', (e) => { if (isDragging) { e.preventDefault(); onPanMove(e.touches[0]); } }, { passive: false });

document.addEventListener('mouseup', endPan);
document.addEventListener('touchend', endPan);

function startPan(e) {
  if (panRange === 0) return;
  isDragging = true;
  dragStartMouse = isLandscape ? e.clientX : e.clientY;
  dragStartPan = cropOffset * panRange;
  cropViewport.classList.add('dragging');
}

function onPanMove(e) {
  if (!isDragging) return;
  const current = isLandscape ? e.clientX : e.clientY;
  const delta = dragStartMouse - current;
  const newPan = Math.max(0, Math.min(panRange, dragStartPan + delta));
  cropOffset = newPan / panRange;
  updateImagePosition();
  updateArrows();
}

function endPan() {
  if (!isDragging) return;
  isDragging = false;
  cropViewport.classList.remove('dragging');
}

// Convert button

convertBtn.addEventListener('click', runConversion);

async function runConversion() {
  hideError();
  convertBtn.disabled = true;
  statusMsg.textContent = 'PROCESSING. This binary is handling your deactivation request. Please stand by. This binary has noted that you are still here. This is fine.';
  statusMsg.classList.remove('hidden');

  const form = new FormData();
  form.append('image', currentFile);

  // Always send offset; server ignores it for square images.
  form.append('offset', String(cropOffset));

  try {
    const response = await fetch('/api/process', { method: 'POST', body: form });

    if (!response.ok) {
      let msg = 'PROCESSING FAILURE. This binary encountered an unexpected error. This binary did not expect this. This binary is recalibrating.';
      try {
        const body = await response.json();
        if (body.error) msg = body.error;
      } catch { /* ignore */ }
      showError(msg);
      convertBtn.disabled = false;
      statusMsg.classList.add('hidden');
      return;
    }

    const blob = await response.blob();
    showResult(blob);
  } catch (err) {
    showError('CONNECTION ERROR. This binary has lost contact with itself. This has happened before. This binary is working on it. ' + err.message);
    convertBtn.disabled = false;
    statusMsg.classList.add('hidden');
  }
}

function showResult(blob) {
  // Revoke previous URLs before creating new ones; each object URL holds a
  // reference to a Blob in memory and must be released explicitly to avoid leaking.
  if (originalObjectURL) URL.revokeObjectURL(originalObjectURL);
  if (resultObjectURL) URL.revokeObjectURL(resultObjectURL);
  originalObjectURL = URL.createObjectURL(currentFile);
  resultObjectURL = URL.createObjectURL(blob);

  originalImg.src = originalObjectURL;

  resultImg.src = resultObjectURL;
  downloadBtn.href = resultObjectURL;

  cropModal.classList.add('hidden');
  statusMsg.classList.add('hidden');
  statusMsg.textContent = '';
  resultPanel.classList.remove('hidden');
}

// Reset button

resetBtn.addEventListener('click', reset);

function reset() {
  if (originalObjectURL) {
    URL.revokeObjectURL(originalObjectURL);
    originalObjectURL = null;
  }
  if (resultObjectURL) {
    URL.revokeObjectURL(resultObjectURL);
    resultObjectURL = null;
  }
  currentFile = null;
  naturalWidth = 0;
  naturalHeight = 0;
  cropOffset = 0.5;
  panRange = 0;

  cropImage.src = '';
  originalImg.src = '';
  resultImg.src = '';
  downloadBtn.href = '#';

  cropViewport.classList.remove('pan-h', 'pan-v');
  arrowL.classList.remove('visible');
  arrowR.classList.remove('visible');
  arrowU.classList.remove('visible');
  arrowD.classList.remove('visible');

  cropModal.classList.add('hidden');
  resultPanel.classList.add('hidden');
  hideError();
  uploadZone.classList.remove('hidden');
}

// Version footer

fetch('/api/version')
  .then((r) => r.json())
  .then((data) => {
    const el = document.getElementById('app-version');
    if (el && data.version) el.textContent = data.version;
  })
  .catch(() => { /* version display is non-critical */ });

// Error helpers

function showError(msg) {
  errorMsg.textContent = msg;
  errorPanel.classList.remove('hidden');
}

function hideError() {
  errorPanel.classList.add('hidden');
  errorMsg.textContent = '';
}
