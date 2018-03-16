window.addEventListener('load', function () {
  document.getElementById('sign-out').onclick = function () {
    firebase.auth().signOut();
  };
  // FirebaseUI config.
  var uiConfig = {
    signInSuccessUrl: '/',
    signInOptions: [
      // Leave the lines as is for the providers you want to offer your users.
      firebase.auth.GoogleAuthProvider.PROVIDER_ID,
      firebase.auth.FacebookAuthProvider.PROVIDER_ID,
      firebase.auth.TwitterAuthProvider.PROVIDER_ID,
      firebase.auth.GithubAuthProvider.PROVIDER_ID,
      firebase.auth.EmailAuthProvider.PROVIDER_ID,
      firebase.auth.PhoneAuthProvider.PROVIDER_ID
    ],
    // Terms of service url.
    tosUrl: '<your-tos-url>'
  };

  firebase.auth().onAuthStateChanged(function (user) {
    if (user) {
      // User is signed in.
      document.getElementById('sign-out').hidden = false;
      document.getElementById('post-form').hidden = false;
      var account = document.getElementById('account-details').textContent =
          'Signed in as ' + user.displayName + ' (' + user.email + ')';
      user.getIdToken().then(function (accessToken) {
        // Add the token to the post form. The user info will be extracted
        // from the token by the server.
        document.getElementById('token').value = accessToken;
      });
    } else {
      // User is signed out.
      // Initialize the FirebaseUI Widget using Firebase.
      var ui = new firebaseui.auth.AuthUI(firebase.auth());
      // Show the Firebase login button.
      ui.start('#firebaseui-auth-container', uiConfig);
      // Update the login state indicators.
      document.getElementById('sign-out').hidden = true;
      document.getElementById('post-form').hidden = true;
      document.getElementById('account-details').textContent = '';
    }
  }, function (error) {
    console.log(error);
    alert('Unable to log in: ' + error)
  });
});