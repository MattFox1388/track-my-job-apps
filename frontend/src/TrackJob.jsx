import { useState } from 'react'
import './TrackJob.css'

function TrackJob() {
  const [jobText, setJobText] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [parsedJob, setParsedJob] = useState(null)
  const [notes, setNotes] = useState('')
  const [isSaving, setIsSaving] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [editedJob, setEditedJob] = useState(null)
  const [platform, setPlatform] = useState('linkedin')
  const [extractedUrl, setExtractedUrl] = useState('')

  // Function to detect mf-URL at bottom of job text and extract the URL
  const detectMfUrl = (text) => {
    const mfUrlMatch = text.match(/mf-URL:\s*([^\s]+)/i)
    
    if (mfUrlMatch) {
      const url = mfUrlMatch[1].trim()
      setExtractedUrl(url)
      
      // Auto-switch to greenhouse if URL contains job-boards.greenhouse.io
      if (url.includes('job-boards.greenhouse.io')) {
        setPlatform('greenhouse')
      }
      if (url.includes('myworkdayjobs.com')) {
        setPlatform('workday')
      }
      
      return url
    }
    
    setExtractedUrl('')
    return null
  }

  const handleJobTextChange = (e) => {
    const newText = e.target.value
    setJobText(newText)
    detectMfUrl(newText)
  }

  const handleTrackJob = async () => {

    if (!jobText.trim()) {
      alert("Please enter job application text")
      return
    }

    setIsLoading(true)
    try {
      // Check if window.go exists
      if (!window.go) {
        throw new Error("Wails runtime not available. Make sure you're running through the Wails app, not the browser.")
      }

      const jobApp = await window.go.main.App.TrackJobApp(jobText.trim(), platform)
      console.log("Parsed job:", jobApp)
      setParsedJob(jobApp)
      setEditedJob(jobApp) // Initialize edited job with parsed data
      setJobText('')
      setNotes(jobApp.notes || '') // Initialize notes with parsed job notes
      setIsEditing(false) // Reset editing state
    } catch (error) {
      console.error("Error calling Go function:", error)
      alert("Error: " + error.message)
    } finally {
      setIsLoading(false)
    }
  }

  const handleSaveJob = async () => {
    setIsSaving(true)
    try {
      // Use edited job data if available, otherwise use parsed job
      const jobToSave = { ...(editedJob || parsedJob), notes: notes.trim() }
      await window.go.main.App.SaveJobApp(jobToSave)
      alert("Job application saved successfully!")
      setParsedJob(null)
      setEditedJob(null)
      setNotes('')
      setIsEditing(false)
    } catch (error) {
      console.error("Error saving job:", error)
      alert("Error saving job: " + error.message)
    } finally {
      setIsSaving(false)
    }
  }

  const handleEditClick = () => {
    setIsEditing(true)
    if (!editedJob) {
      setEditedJob({ ...parsedJob })
    }
  }

  const handleCancelEdit = () => {
    setIsEditing(false)
    setEditedJob({ ...parsedJob }) // Reset to original parsed data
  }

  const handleFieldChange = (field, value) => {
    setEditedJob(prev => ({
      ...prev,
      [field]: value
    }))
  }

  return (
    <div className="container">
      <h1>Track My Job Apps</h1>

      <div className="input-container">
        <div style={{ marginBottom: '15px', display: 'flex', alignItems: 'center', gap: '15px', justifyContent: 'center' }}>
          <label style={{ fontWeight: '500', color: 'white', fontSize: '16px' }}>Platform:</label>
          <select
            value={platform}
            onChange={(e) => setPlatform(e.target.value)}
            disabled={isLoading}
            style={{
              padding: '8px 12px',
              borderRadius: '8px',
              border: '1px solid rgba(255, 255, 255, 0.3)',
              background: 'rgba(255, 255, 255, 0.1)',
              color: 'white',
              fontSize: '14px',
              cursor: 'pointer',
              minWidth: '120px'
            }}
          >
            <option value="linkedin" style={{ background: '#333', color: 'white' }}>LinkedIn</option>
            <option value="greenhouse" style={{ background: '#333', color: 'white' }}>Greenhouse</option>
            <option value="workday" style={{ background: '#333', color: 'white' }}>Workday</option>
            <option value="other" style={{ background: '#333', color: 'white' }}>Other</option>
          </select>
        </div>
        
        <textarea
          value={jobText}
          onChange={handleJobTextChange}
          placeholder={`Enter a job application from ${platform === 'linkedin' ? 'LinkedIn' : platform === 'greenhouse' ? 'Greenhouse' : 'Other'}...`}
          disabled={isLoading}
        />
        {extractedUrl && (
          <div style={{ 
            marginTop: '10px', 
            padding: '10px', 
            background: 'rgba(76, 175, 80, 0.2)', 
            borderRadius: '5px',
            border: '1px solid rgba(76, 175, 80, 0.5)',
            color: '#4CAF50'
          }}>
            <strong>✓ URL detected:</strong> {extractedUrl}
            {platform === 'greenhouse' && (
              <div style={{ marginTop: '5px', fontSize: '14px' }}>
                <em>Platform automatically switched to Greenhouse</em>
              </div>
            )}
          </div>
        )}
        <button
          onClick={handleTrackJob}
          disabled={isLoading || !jobText.trim()}
        >
          {isLoading ? 'Parsing...' : 'Parse Job'}
        </button>
      </div>

      {parsedJob && (
        <>
          <div style={{ marginTop: '20px', textAlign: 'left', background: 'rgba(255,255,255,0.1)', padding: '15px', borderRadius: '10px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
              <h3 style={{ margin: 0 }}>Job Application Details:</h3>
              {!isEditing ? (
                <button 
                  onClick={handleEditClick}
                  style={{
                    backgroundColor: 'rgba(255,255,255,0.2)',
                    color: 'white',
                    border: '1px solid rgba(255,255,255,0.3)',
                    borderRadius: '5px',
                    padding: '8px 16px',
                    cursor: 'pointer',
                    fontSize: '14px'
                  }}
                >
                  ✏️ Edit
                </button>
              ) : (
                <div style={{ display: 'flex', gap: '10px' }}>
                  <button 
                    onClick={handleCancelEdit}
                    style={{
                      backgroundColor: 'rgba(255,255,255,0.2)',
                      color: 'white',
                      border: '1px solid rgba(255,255,255,0.3)',
                      borderRadius: '5px',
                      padding: '8px 16px',
                      cursor: 'pointer',
                      fontSize: '14px'
                    }}
                  >
                    ❌ Cancel
                  </button>
                </div>
              )}
            </div>

            {isEditing ? (
              <div style={{ display: 'grid', gap: '15px' }}>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>Company:</label>
                  <input
                    type="text"
                    value={editedJob?.company || ''}
                    onChange={(e) => handleFieldChange('company', e.target.value)}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      borderRadius: '5px',
                      border: '1px solid rgba(255,255,255,0.3)',
                      background: 'rgba(255,255,255,0.1)',
                      color: 'white',
                      fontSize: '14px'
                    }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>Position:</label>
                  <input
                    type="text"
                    value={editedJob?.position || ''}
                    onChange={(e) => handleFieldChange('position', e.target.value)}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      borderRadius: '5px',
                      border: '1px solid rgba(255,255,255,0.3)',
                      background: 'rgba(255,255,255,0.1)',
                      color: 'white',
                      fontSize: '14px'
                    }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>Location:</label>
                  <input
                    type="text"
                    value={editedJob?.location || ''}
                    onChange={(e) => handleFieldChange('location', e.target.value)}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      borderRadius: '5px',
                      border: '1px solid rgba(255,255,255,0.3)',
                      background: 'rgba(255,255,255,0.1)',
                      color: 'white',
                      fontSize: '14px'
                    }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>Salary Range:</label>
                  <input
                    type="text"
                    value={editedJob?.salaryRange || ''}
                    onChange={(e) => handleFieldChange('salaryRange', e.target.value)}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      borderRadius: '5px',
                      border: '1px solid rgba(255,255,255,0.3)',
                      background: 'rgba(255,255,255,0.1)',
                      color: 'white',
                      fontSize: '14px'
                    }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>Workplace Type:</label>
                  <input
                    type="text"
                    value={editedJob?.workplaceType || ''}
                    onChange={(e) => handleFieldChange('workplaceType', e.target.value)}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      borderRadius: '5px',
                      border: '1px solid rgba(255,255,255,0.3)',
                      background: 'rgba(255,255,255,0.1)',
                      color: 'white',
                      fontSize: '14px'
                    }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>Status:</label>
                  <select
                    value={editedJob?.status || 'SUBMITTED'}
                    onChange={(e) => handleFieldChange('status', e.target.value)}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      borderRadius: '5px',
                      border: '1px solid rgba(255,255,255,0.3)',
                      background: 'rgba(255,255,255,0.1)',
                      color: 'white',
                      fontSize: '14px'
                    }}
                  >
                    <option value="SUBMITTED">SUBMITTED</option>
                    <option value="REJECTED">REJECTED</option>
                    <option value="PHONE_SCREEN">PHONE_SCREEN</option>
                    <option value="REMOTE_INTERVIEW">REMOTE_INTERVIEW</option>
                    <option value="ON_SITE_INTERVIEW">ON_SITE_INTERVIEW</option>
                  </select>
                </div>
              </div>
            ) : (
              <div>
                <p><strong>Company:</strong> {(editedJob || parsedJob).company || 'Not found'}</p>
                <p><strong>Position:</strong> {(editedJob || parsedJob).position || 'Not found'}</p>
                <p><strong>Location:</strong> {(editedJob || parsedJob).location || 'Not found'}</p>
                <p><strong>Salary Range:</strong> {(editedJob || parsedJob).salaryRange || 'Not found'}</p>
                <p><strong>Workplace Type:</strong> {(editedJob || parsedJob).workplaceType || 'Not found'}</p>
                <p><strong>Status:</strong> {(editedJob || parsedJob).status}</p>
              </div>
            )}
          </div>
          
          <div style={{ marginTop: '20px' }}>
            <h4>Add Notes (Optional):</h4>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Add any additional notes about this job application..."
              style={{
                width: '100%',
                minHeight: '80px',
                padding: '10px',
                borderRadius: '5px',
                border: '1px solid #ccc',
                fontSize: '14px',
                marginBottom: '15px',
                resize: 'vertical'
              }}
              disabled={isSaving}
            />
            <button 
              onClick={handleSaveJob}
              disabled={isSaving}
              style={{
                backgroundColor: isSaving ? '#666' : '#4CAF50',
                color: 'white',
                padding: '10px 20px',
                border: 'none',
                borderRadius: '5px',
                cursor: isSaving ? 'not-allowed' : 'pointer',
                fontSize: '16px'
              }}
            >
              {isSaving ? 'Saving...' : 'Save Job Application'}
            </button>
          </div>
        </>
      )}
    </div>
  )
}

export default TrackJob