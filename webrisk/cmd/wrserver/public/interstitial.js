// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/**
 * Alias for document.getElementById. Found elements must be HTMLElements.
 * @param {string} id The ID of the element to find.
 * @return {HTMLElement} The found element or null if not found.
 */
function $(id) {
  var el = document.getElementById(id);
  return el;
}

/** Namespace. */
var interstitial = interstitial || {};


/**
 * Whether the page is currently being viewed at a "mobile" screen size.
 * @type {boolean}
 * @private
 */
interstitial.mobileNav_ = false;


/**
 * Set up event handlers for UI elements.
 */
interstitial.setupEvents = function() {
  // The "back to safety" button.
  $('primary-button').addEventListener('click', function() {
    window.history.back();
  });

  // The "Details" button.
  $('details-button').addEventListener('click', function(event) {
    var hiddenDetails = $('details').classList.toggle('hidden');

    if (interstitial.mobileNav_)
      $('main-content').classList.toggle('hidden', !hiddenDetails);
    else
      $('main-content').classList.remove('hidden');

    $('details-button').innerText = hiddenDetails ? 'Details' : 'Hide details';
  });

  // Handle resize events.
  window.addEventListener('resize', interstitial.onResize);
  interstitial.onResize();
};


/**
 * For small screen mobile, the navigation buttons are moved below the advanced
 * text.
 */
interstitial.onResize = function() {
  var helpOuterBox = document.querySelector('#details');
  var mainContent = document.querySelector('#main-content');
  var mediaQuery = '(min-width: 240px) and (max-width: 420px) and ' +
      '(max-height: 736px) and (min-height: 401px) and ' +
      '(orientation: portrait), (max-width: 736px) and ' +
      '(max-height: 420px) and (min-height: 240px) and ' +
      '(min-width: 421px) and (orientation: landscape)';

  var detailsHidden = helpOuterBox.classList.contains('hidden');

  // Check for change in nav status.
  if (interstitial.mobileNav_ != window.matchMedia(mediaQuery).matches) {
    interstitial.mobileNav_ = !interstitial.mobileNav_;

    if (interstitial.mobileNav_) {
      mainContent.classList.toggle('hidden', !detailsHidden);
      helpOuterBox.classList.toggle('hidden', detailsHidden);
    } else if (!detailsHidden) {
      mainContent.classList.remove('hidden');
      helpOuterBox.classList.remove('hidden');
    }
  }
};


document.addEventListener('DOMContentLoaded', interstitial.setupEvents);
