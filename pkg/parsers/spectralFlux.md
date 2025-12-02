The most well-known and widely used algorithm for generating an Onset Strength Envelope is **Spectral Flux** (specifically, **Log-Spectral Flux**).

This algorithm was popularized in music information retrieval (MIR) by researchers like **Juan Bello** and is the default method used in libraries like Python’s **Librosa**.

### Why not just use volume?
You might think you could just look at the volume (amplitude) of the waveform to find beats. The problem is that a long, loud synthesizer chord has high volume, but no "beat." The Onset Envelope needs to capture **sudden changes**, not just loudness.

---

### The Spectral Flux Algorithm (Step-by-Step)

This algorithm transforms the audio based on the idea that **a beat occurs when energy suddenly increases across many frequencies.**

#### 1. Short-Time Fourier Transform (STFT)
First, the algorithm breaks the waveform into tiny overlapping windows (usually ~20-30 milliseconds) and calculates the frequency spectrum for each window.
*   **Input:** Time-domain signal (Amplitude over time).
*   **Output:** Spectrogram (Frequency magnitude over time).

#### 2. Logarithmic Compression
Human hearing is logarithmic (we hear volume in decibels, not linear amplitude). The algorithm converts the linear magnitudes from the STFT into a logarithmic scale. This makes weak beats (like hi-hats) detectable even if a loud bass line is playing.
*   **Formula:** $Y = \log(1 + \text{magnitude})$

#### 3. Differentiation (The "Flux")
This is the core of the algorithm. It compares the spectrum of the **current frame** with the spectrum of the **previous frame**.
It calculates the difference bin-by-bin.
*   *If 100Hz was quiet in frame 1, and loud in frame 2 = **High Flux**.*
*   *If 100Hz was loud in frame 1, and stays loud in frame 2 = **Zero Flux** (Sustain).*

#### 4. Half-Wave Rectification
We are only interested in the *start* of a sound (energy increase), not the end (energy decrease).
Therefore, the algorithm throws away any negative values. If the energy dropped, the result is 0. If the energy rose, we keep the value.

#### 5. Summation
Finally, it sums up the positive changes across all frequency bins for that specific moment in time. This single number represents the "Onset Strength" for that frame.

---

### Summary Formula
If you were to write this mathematically, the Onset Strength $O[n]$ at time frame $n$ looks like this:

$$ O[n] = \sum_{k=0}^{K} \max(0, \ Y(n, k) - Y(n-1, k)) $$

Where:
*   $Y(n, k)$ is the Log-Magnitude of frequency bin $k$ at time $n$.
*   The $\max(0, ...)$ part ensures we only count energy **increases**.

### A Modern Improvement: "Superflux"
While standard Spectral Flux is the classic method, a modern variation called **Superflux** (developed by Böck and Widmer) is now often preferred (and available in Librosa).

It addresses a major flaw in standard Spectral Flux: **Vibrato.**
If a singer holds a note but shakes their pitch (vibrato), the energy moves from one frequency bin to another. Standard Spectral Flux interprets this as "Old frequency went down, new frequency went up" and mistakenly triggers an onset.

**Superflux** applies a "maximum filter" (it looks at adjacent frequency bands) before calculating the difference, effectively suppressing vibrato and pitch sliding so that only *true* percussive hits are registered.