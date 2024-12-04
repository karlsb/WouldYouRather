/** @type {import('tailwindcss').Config} */
import daisyui from 'daisyui'
import tailwindcssanimate from 'tailwindcss-animate' 
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'text': '#FF6E31',
        'secondary': '#243763',
        'primary': '#FFEBB7',
        'tertiary': '#AD8E70', 
      },
      scale: {
        '120': '1.2',
      },
      animation: {
				fade: 'fadeIn .3s ease-in-out',
        enlarge: 'large .3s ease-in',
			},
			keyframes: {
				fadeIn: {
					from: { opacity: 0 },
					to: { opacity: 1 },
				},
        large: {
          from: {transform: 'scale(1.0)'},
          to:{transform:'scale(1.2)'}
        }
			},
    },
  },
  plugins: [
    daisyui,
    tailwindcssanimate,
  ],
}

