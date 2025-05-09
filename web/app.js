document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('show-subjects-btn').addEventListener('click', async () => {
        try {
            const res = await fetch('/api/subjects/devices');
            const data = await res.json();
            if (!data.status) {
                throw new Error('Gagal mengambil data dari database');
            }

            const list = document.getElementById('subject-list');
            list.innerHTML = '';
            data.data.forEach(subject => {
                const s = document.createElement('div');
                s.classList.add('subject-row');
                
                const device = document.createElement('div');
                device.textContent = subject.device_id;
                device.classList.add('subject-data');

                const name = document.createElement('div');
                name.textContent = subject.name;
                name.classList.add('subject-data');
                
                const isFatigued = document.createElement('div');
                isFatigued.textContent = subject.is_fatigued;
                isFatigued.classList.add('subject-data');

                const createdAt = document.createElement('div');
                createdAt.textContent = new Date(subject.created_at).toLocaleString();
                createdAt.classList.add('subject-data');

                s.appendChild(device);
                s.appendChild(name);
                s.appendChild(isFatigued);
                s.appendChild(createdAt);

                list.appendChild(s);
            });

            document.getElementById('modal').classList.remove('hidden');
        } catch (err) {
            console.error('Gagal memuat data: ', err);
        }
    });

    document.getElementById('close-modal').addEventListener('click', function(){
        document.getElementById('modal').classList.add('hidden');
    });
});