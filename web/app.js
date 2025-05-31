const setupRealTimeUpdates = () => {
    // Pilihan 1: Polling
    setInterval(async () => {
      try {
        const res = await fetch('/api/devices/subjects');
        const data = await res.json();
        
        // Perbarui tampilan jika ada perubahan
        const currentCount = document.querySelectorAll('.device-subject-row-container').length;
        if (data.data.length !== currentCount) {
          window.location.reload(); // Atau update bagian tertentu saja
        }
      } catch (err) {
        console.error('Error checking updates:', err);
      }
    }, 4000); 
  };

document.addEventListener('DOMContentLoaded', async () => {
    setupRealTimeUpdates();
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

                const deleteBtn = document.createElement('button');
                deleteBtn.textContent = 'X';
                deleteBtn.classList.add('delete-subject-btn');
                deleteBtn.title = 'Hapus subject ini';
                
                deleteBtn.addEventListener('click', async () => {
                    const confirmed = confirm(`Yakin ingin menghapus subyek "${subject.name}"?`);
                    if (!confirmed) return;
            
                    try {
                        const res = await fetch(`/api/subjects/${subject.id}`, {
                            method: 'DELETE'
                        });
                        const result = await res.json();
            
                        if (result.status) {
                            alert('Subyek berhasil dihapus!');
                            s.remove(); // Hapus dari list modal
            
                            const options = document.querySelectorAll('#set-subject option');
                            options.forEach(opt => {
                                if (parseInt(opt.value) === subject.id) {
                                    opt.remove();
                                }
                            });
            
                        } else {
                            alert('Gagal menghapus subyek: ' + result.message);
                        }
                    } catch (err) {
                        console.error('Gagal menghapus subyek: ', err);
                        alert('Terjadi kesalahan saat menghapus subyek.');
                    }
                });

                s.appendChild(name);
                s.appendChild(createdAt);
                s.appendChild(deleteBtn);

                list.appendChild(s);
            });

            document.getElementById('subject-modal').classList.remove('hidden');
        } catch (err) {
            console.error('Gagal memuat data: ', err);
        }
    });

    document.getElementById('close-subject-modal').addEventListener('click', function () {
        document.getElementById('subject-modal').classList.add('hidden');
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

    const setModal = document.getElementById('set-modal');
    const select = document.getElementById('set-subject');
    fetch('/api/subjects')
    .then(response => response.json())
    .then(data => {
        data.data.forEach(subject => {
            const opt = document.createElement('option');
            opt.value = subject.id;
            opt.textContent = subject.name;
            select.appendChild(opt);
        })
    });

    try {
        const resDs = await fetch('/api/devices/subjects');
        const dataDs = await resDs.json();

        const dsList = document.getElementById('device-subject-list');
        dsList.innerHTML = '';
        dataDs.data.forEach(item => {
            const row = document.createElement('div');
            row.classList.add('device-subject-row-container');
            const rowData = document.createElement('div');
            rowData.classList.add('device-subject-row');

            const noSubject = item.subject_id === 0;
            const dId = item.device_id;

            const deviceId = document.createElement('div');
            deviceId.textContent = dId;
            deviceId.classList.add('device-subject-data');
            
            const status = document.createElement('div');
            status.textContent = item.device_status;
            status.classList.add('device-subject-data', 'device-status');

            const name = document.createElement('div');
            name.textContent = noSubject ? '' : item.name;
            name.classList.add('device-subject-data');

            const isFatigued = document.createElement('div');
            isFatigued.textContent = noSubject ? '' : (item.is_fatigued ? "lelah" : "tidak lelah"); // Konversi ke teks
            isFatigued.classList.add('device-subject-data');
            isFatigued.setAttribute('data-subject-status', item.is_fatigued.toString()); // Tambahkan atribut data

            const createdAt = document.createElement('div');
            createdAt.textContent = noSubject ? '' : new Date(item.created_at).toLocaleString();
            createdAt.classList.add('device-subject-data');

            const actionButtenContainer = document.createElement('div');
            actionButtenContainer.classList.add('action-btn-container');

            const setSubjectButton = document.createElement('button');
            setSubjectButton.textContent = 'Atur Subyek';
            setSubjectButton.classList.add('set-subject-btn');

            setSubjectButton.addEventListener('click', (e) => {
                const rect = e.target.getBoundingClientRect();
                
                setModal.style.position = 'absolute';
                setModal.style.top = `${rect.bottom + window.scrollY}px`;
                setModal.style.left = `${rect.left + window.scrollX}px`;
                
                setModal.classList.toggle('hidden');

                const existingBtn = document.getElementById('submit-set-btn');
                if (existingBtn) existingBtn.remove();

                const setBtn = document.createElement('button');
                setBtn.textContent = 'Simpan';
                setBtn.id = 'submit-set-btn';
                setBtn.addEventListener('click', (e) => {
                    e.preventDefault();

                    fetch('/api/devices/' + dId + '/subjects', {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            subject_id: parseInt(select.value)
                        })
                    })
                    .then(res => res.json())
                    .then(data => {
                        select.value = '';
                        setModal.classList.add('hidden');
                        if (data.status) {
                            alert('Subyek berhasil dipasangkan');
                            window.location.reload();
                        } else {
                            alert('Gagal memasangkan subyek..' + data.message);
                        }
                    })
                    .catch(err => {
                        console.error('Fetch error: ', err);
                    });
                });

                setModal.appendChild(setBtn);
            });

            const deleteBtn = document.createElement('button');
            deleteBtn.textContent = 'hapus perangkat'; // atau pakai icon lain
            deleteBtn.classList.add('delete-device-btn');
            deleteBtn.title = 'Hapus perangkat';

            // Event untuk menghapus device
            deleteBtn.addEventListener('click', async () => {
                const confirmed = confirm(`Yakin ingin menghapus device ${dId}?`);
                if (!confirmed) return;

                try {
                    const res = await fetch(`/api/devices/${dId}`, {
                        method: 'DELETE'
                    });
                    const result = await res.json();

                    if (result.status) {
                        alert('Device berhasil dihapus!');
                        row.remove(); // Hapus dari tampilan
                    } else {
                        alert('Gagal menghapus device: ' + result.message);
                    }
                } catch (err) {
                    console.error('Gagal menghapus device: ', err);
                    alert('Terjadi kesalahan saat menghapus device.');
                }
            });

            rowData.appendChild(deviceId);
            rowData.appendChild(status);
            rowData.appendChild(name);
            rowData.appendChild(isFatigued);
            rowData.appendChild(createdAt);
            actionButtenContainer.appendChild(setSubjectButton);
            actionButtenContainer.appendChild(deleteBtn);
            
            row.appendChild(rowData);
            row.appendChild(actionButtenContainer);
   


            dsList.appendChild(row);
        });
    } catch (err) {
        console.error('Gagal memuat device-subjects: ', err);
    }

    document.getElementById('refresh-btn').addEventListener('click', () => {
        window.location.reload();
    });

    
    
});
