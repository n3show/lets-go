(function () {
  var toggle = document.getElementById('navToggle');
  var sidebar = document.getElementById('sidebar');
  var scrim = document.getElementById('scrim');

  function closeNav() {
    sidebar.classList.remove('is-open');
    scrim.classList.remove('is-open');
    toggle.setAttribute('aria-expanded', 'false');
  }

  function openNav() {
    sidebar.classList.add('is-open');
    scrim.classList.add('is-open');
    toggle.setAttribute('aria-expanded', 'true');
  }

  if (toggle) {
    toggle.addEventListener('click', function () {
      var isOpen = sidebar.classList.contains('is-open');
      if (isOpen) { closeNav(); } else { openNav(); }
    });
  }
  if (scrim) {
    scrim.addEventListener('click', closeNav);
  }

  // close drawer when a nav link is tapped
  document.querySelectorAll('.toc a').forEach(function (a) {
    a.addEventListener('click', closeNav);
  });

  // copy buttons
  document.querySelectorAll('.copy-btn').forEach(function (btn) {
    btn.addEventListener('click', function () {
      var code = btn.parentElement.querySelector('code');
      var text = code ? code.innerText : '';
      navigator.clipboard.writeText(text).then(function () {
        var original = btn.textContent;
        btn.textContent = 'Copied';
        setTimeout(function () { btn.textContent = original; }, 1200);
      });
    });
  });

  // theme toggle
  var themeToggle = document.getElementById('themeToggle');
  if (themeToggle) {
    themeToggle.addEventListener('click', function () {
      var root = document.documentElement;
      var current = root.getAttribute('data-theme') === 'dark' ? 'dark' : 'light';
      var next = current === 'dark' ? 'light' : 'dark';
      root.setAttribute('data-theme', next);
      localStorage.setItem('lg-theme', next);
    });
  }
})();
