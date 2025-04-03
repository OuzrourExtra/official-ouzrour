function updateTimeRemaining() {
    document.querySelectorAll("[data-id]").forEach(el => {
        const id = el.dataset.id;
        const created = new Date(el.dataset.created);
        const planned = new Date(el.dataset.planned);
        const now = new Date();

        const totalMs = planned.getTime() - created.getTime();
        const remainingMs = planned.getTime() - now.getTime();

        const progress = Math.max(0, Math.min(100, Math.floor((remainingMs / totalMs) * 100)));

        const bar = document.getElementById(`timeProgressBar-${id}`);
        const label = document.getElementById(`timeProgressText-${id}`);

        // Affichage de la barre (progression du temps restant)
        if (bar) bar.style.width = `${progress}%`;

        // Si la date est dépassée
        if (remainingMs <= 0) {
            if (label) label.textContent = "⏰ Échéance dépassée";
            if (bar) bar.classList.add("bg-red-600");
            return;
        }

        // Formatage du temps restant
        let remaining = remainingMs;

        const days = Math.floor(remaining / (1000 * 60 * 60 * 24));
        remaining -= days * 24 * 60 * 60 * 1000;

        const hours = Math.floor(remaining / (1000 * 60 * 60));
        remaining -= hours * 60 * 60 * 1000;

        const minutes = Math.floor(remaining / (1000 * 60));
        remaining -= minutes * 60 * 1000;

        const seconds = Math.floor(remaining / 1000);

        const pad = n => String(n).padStart(2, "0");

        if (label) {
            label.textContent = `${pad(days)}j : ${pad(hours)}h : ${pad(minutes)}min : ${pad(seconds)}s`;
        }
    });
}

updateTimeRemaining();
setInterval(updateTimeRemaining, 1000);