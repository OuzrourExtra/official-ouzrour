function updateTimeRemaining() {
    document.querySelectorAll("[data-id]").forEach(el => {
        const id = el.dataset.id;
        const created = new Date(el.dataset.created);
        const planned = new Date(el.dataset.planned);
        const now = new Date();

        const totalMs = planned.getTime() - created.getTime();
        const elapsedMs = now.getTime() - created.getTime();
        const remainingMs = planned.getTime() - now.getTime();

        // Calcul du pourcentage de temps ÉCOULÉ (et non temps restant)
        let progress = Math.floor((elapsedMs / totalMs) * 100);
        progress = Math.max(0, Math.min(100, progress)); // Clamp 0-100

        const bar = document.getElementById(`timeProgressBar-${id}`);
        const label = document.getElementById(`timeProgressText-${id}`);

        // Barre de progression du TEMPS ÉCOULÉ
        if (bar) bar.style.width = `${progress}%`;

        // Si date dépassée
        if (remainingMs <= 0) {
            if (label) label.textContent = "⏰ Échéance dépassée";
            if (bar) bar.classList.add("bg-red-600");
            return;
        }

        // Affichage temps restant (en texte)
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
