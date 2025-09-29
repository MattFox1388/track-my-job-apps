import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import MenuItem from '@mui/material/MenuItem';
import { useState, useEffect } from 'react';

const statusOptions = [
    { value: 'SUBMITTED', label: 'Submitted' },
    { value: 'REJECTED', label: 'Rejected' },
    { value: 'PHONE_SCREEN', label: 'Phone Screen' },
    { value: 'REMOTE_INTERVIEW', label: 'Remote Interview' },
    { value: 'ON_SITE_INTERVIEW', label: 'On-Site Interview' }
];

export default function EditModal(props) {
    const [formData, setFormData] = useState({
        appId: '',
        company: '',
        position: '',
        location: '',
        salaryRange: '',
        workplaceType: '',
        status: '',
        notes: ''
    });

    useEffect(() => {
        if (props.jobApp) {
            setFormData({
                appId: props.jobApp.appId || '',
                company: props.jobApp.company || '',
                position: props.jobApp.position || '',
                location: props.jobApp.location || '',
                salaryRange: props.jobApp.salaryRange || '',
                workplaceType: props.jobApp.workplaceType || '',
                status: props.jobApp.status || '',
                notes: props.jobApp.notes || ''
            });
        }
    }, [props.jobApp]);

    const handleChange = (field) => (event) => {
        setFormData(prev => ({
            ...prev,
            [field]: event.target.value
        }));
    };

    const handleSubmit = () => {
        if (props.onSave) {
            props.onSave(formData);
        }
        props.handleClose();
    };

    return (
        <Dialog open={props.open} onClose={props.handleClose} maxWidth="sm" fullWidth>
            <DialogTitle>Edit Job Application</DialogTitle>
            <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
                <TextField 
                    label="Company" 
                    value={formData.company}
                    onChange={handleChange('company')}
                    fullWidth
                    variant="outlined"
                    sx={{ marginTop: '10px' }}
                />
                <TextField 
                    label="Position" 
                    value={formData.position}
                    onChange={handleChange('position')}
                    fullWidth
                    variant="outlined"
                />
                <TextField 
                    label="Location" 
                    value={formData.location}
                    onChange={handleChange('location')}
                    fullWidth
                    variant="outlined"
                />
                <TextField 
                    label="Salary Range" 
                    value={formData.salaryRange}
                    onChange={handleChange('salaryRange')}
                    fullWidth
                    variant="outlined"
                />
                <TextField 
                    label="Workplace Type" 
                    value={formData.workplaceType}
                    onChange={handleChange('workplaceType')}
                    fullWidth
                    variant="outlined"
                />
                <TextField 
                    label="Status" 
                    value={formData.status}
                    onChange={handleChange('status')}
                    fullWidth
                    variant="outlined"
                    select
                >
                    {statusOptions.map((option) => (
                        <MenuItem key={option.value} value={option.value}>
                            {option.label}
                        </MenuItem>
                    ))}
                </TextField>
                <TextField 
                    label="Notes" 
                    value={formData.notes}
                    onChange={handleChange('notes')}
                    fullWidth
                    multiline
                    rows={3}
                    variant="outlined"
                />
            </DialogContent>
            <DialogActions sx={{ p: 2, gap: 1 }}>
                <Button onClick={props.handleClose} color="inherit">Cancel</Button>
                <Button onClick={handleSubmit} variant="contained" color="primary">Save</Button>
            </DialogActions>
        </Dialog>
    )
}