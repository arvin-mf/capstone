document.addEventListener('DOMContentLoaded', async () => {
    document.getElementById('show-subjects-btn').addEventListener('click', async () => {
        try {
            const res = await fetch('/api/subjects');
            const data = await res.json();
            if (!data.status) {
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

                const createdAt = document.createElement('div');
                createdAt.textContent = new Date(subject.created_at).toLocaleString();
                createdAt.classList.add('subject-data');

                s.appendChild(name);
                s.appendChild(createdAt);

                list.appendChild(s);
            });

            document.getElementById('modal').classList.remove('hidden');
        } catch (err) {
            console.error('Gagal memuat data: ', err);
        }
    });

    document.getElementById('close-modal').addEventListener('click', function () {
        document.getElementById('modal').classList.add('hidden');
    });

    document.getElementById('add-subject-btn').addEventListener('click', () => {
        document.getElementById('subject-form-container').classList.toggle('hidden');
    });

    document.getElementById('submit-subject-btn').addEventListener('click', () => {
        const name = document.getElementById('subject-name').value.trim();

        if (!name) {
            alert('Nama tidak boleh kosong!');
            return;
        }

        fetch('/api/subjects', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        })
        .then(res => res.json())
        .then(data => {
            document.getElementById('subject-name').value = '';
            document.getElementById('add-subject-form').classList.add('hidden');
            if (data.status) {
                alert('Subyek ditambahkan!');
            } else {
                alert('Gagal menambahkan subyek..' + data.message);
            }
        })
        .catch(err => {
            console.error('Fetch error: ', err);
        });
    });

    try {
        const res_ds = await fetch('/api/devices/subjects');
        const data_ds = await res_ds.json();

        const ds_list = document.getElementById('device-subject-list');
        ds_list.innerHTML = '';
        data_ds.data.forEach(item => {
            const row = document.createElement('div');
            row.classList.add('device-subject-row');

            const noSubject = item.subject_id === 0;
            if (noSubject) {
                row.classList.add('no-subject');
            }

            const deviceId = document.createElement('div');
            deviceId.textContent = item.device_id;
            deviceId.classList.add('device-subject-data');

            const name = document.createElement('div');
            name.textContent = noSubject ? '' : item.name;
            name.classList.add('device-subject-data');

            const isFatigued = document.createElement('div');
            isFatigued.textContent = noSubject ? '' : item.is_fatigued;
            isFatigued.classList.add('device-subject-data');

            const createdAt = document.createElement('div');
            createdAt.textContent = noSubject ? '' : new Date(item.created_at).toLocaleString();
            createdAt.classList.add('device-subject-data');

            row.appendChild(deviceId);
            row.appendChild(name);
            row.appendChild(isFatigued);
            row.appendChild(createdAt);

            ds_list.appendChild(row);
        });
    } catch (err) {
        console.error('Gagal memuat device-subjects: ', err);
    }
});
