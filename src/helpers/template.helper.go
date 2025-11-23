package helpers

type Request struct {
	Email string `json:"email"`
}

func Template(req Request) string {
	return `<!DOCTYPE html>
	<html>
	<head>
	<title>Fresh Groove Talent Search</title>
	</head>
	<body>
	<h2>Notification</h2>
	<img src="https://i.ibb.co/1r82KLS/groove.png" alt="Magic 89.9" width="450" height="100" style="margin-right: 20px;">
	<p>Dear ` + req.Email + `</p>
	<p>Thank you for your interest in joining The Fresh Groove online talent search, brought to you by Magic 89.9 and Everything Entertainment. We are on the lookout for the freshest R&B and Pop artists, aged 16 to 22.</p>
	<p><strong>Submit Your Profile Picture</strong> – Upload a high-quality, clear image of yourself. This will be used for your profile in the contest.</p>
	<p><strong>Share Your Talent</strong> – Record a short video (maximum 60 seconds) showcasing your R&B or Pop performance. Make sure it highlights your unique style!</p>
	<p><strong>Fill Out the Registration Form</strong> – Complete the form with your contact details and other required information.</p>
	<ul>
		<li>Name:</li>
		<li>Address:</li>
		<li>Contact Number:</li>
		<li>Date of Birth:</li>
		<li>School:</li>
		<li>Work:</li>
	</ul>
	<p><strong>Submit Your Application</strong> – Attach 3 cover songs in MP3 format and once everything is ready, email your requirements to <a href="mailto:magicfreshgroove@gmail.com">magicfreshgroove@gmail.com</a>, and you're officially entered into the competition!</p>
	<p><strong>Judging Criteria</strong></p>
	<p>The judges will assess each contestant based on the following criteria:</p>
	<ul>
		<li><strong>Vocal Ability (30%)</strong>
			<ul>
				<li>Tone quality: Clear, rich, and pleasant voice.</li>
				<li>Pitch accuracy: Consistent and precise pitch control throughout the performance.</li>
				<li>Vocal range: Ability to sing across various vocal registers (low to high notes).</li>
				<li>Technique: Use of proper breathing, phrasing, and dynamics.</li>
			</ul>
		</li>
		<li><strong>Performance & Stage Presence (25%)</strong>
			<ul>
				<li>Engagement: Ability to connect with the audience, both visually and emotionally.</li>
				<li>Confidence: Displaying self-assurance and comfort on camera.</li>
				<li>Expression: Showing emotion and passion in the performance.</li>
				<li>Movement: Creative and natural gestures or movement to complement the singing.</li>
			</ul>
		</li>
		<li><strong>Song Interpretation (20%)</strong>
			<ul>
				<li>Expression of the song's message: How well the performer conveys the meaning of the lyrics.</li>
				<li>Creativity: Personal interpretation and unique touch to the song selection.</li>
				<li>Adaptation: How well the contestant makes the song their own, while respecting the original composition.</li>
			</ul>
		</li>
		<li><strong>Originality (15%)</strong>
			<ul>
				<li>Unique style: How the performer brings their own personality and flair to the performance.</li>
				<li>Arrangement: Originality in vocal arrangement, if applicable.</li>
			</ul>
		</li>
		<li><strong>Online Voting (10%)</strong>
			<ul>
				<li>Public Appeal: The number of votes a contestant receives from online audiences.</li>
				<li>Engagement: How well the contestant interacts with their supporters and encourages voting.</li>
			</ul>
		</li>
	</ul>
	<p>Good luck,</p>
	<p>The Fresh Groove Team</p>
	<div style="display: flex; justify-content: space-between; gap: 10px; align-items: center;">
		<img src="https://i.ibb.co/T1WB2Mz/everything.png" alt="Entertainment" width="100" height="100" style="margin-right: 20px;">
		<img src="https://i.ibb.co/Ln1yGRk/viber-image-2025-01-12-18-43-52-193.png" alt="Magic Studio Powered by. MFORE" width="200" height="100" style="margin-right: 20px;">
	</div>
</body>
</html>
`
}

// Template function generates an HTML email template for the Fresh Groove Talent Search Auto Reply.
