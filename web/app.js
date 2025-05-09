document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('show-subjects-btn').addEventListener('click', async () => {
        try {
            const res = await fetch('/api/subjects');
            const data = await res.json();
            if (!data.success) {
                throw new Error('Gagal mengambil data dari database');
            }

            const list = document.getElementById('subject-list');
            list.innerHTML = '';
            data.data.forEach(subject => {
                const s = document.createElement('div');
                s.classList.add('subject-row');
                
                const name = document.createElement('div');
                name.textContent = subject.name;
                name.classList.add('subject-data');

                const is_fatigued = document.createElement('div');
                is_fatigued.textContent = subject.is_fatigued;
                is_fatigued.classList.add('subject-data');

                s.appendChild(name);
                s.appendChild(is_fatigued);

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