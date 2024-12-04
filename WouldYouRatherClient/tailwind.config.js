/** @type {import('tailwindcss').Config} */
import daisyui from 'daisyui'
import tailwindcssanimate from 'tailwindcss-animate' 
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  daisyui: {
    themes:[
      {
        one: {
          primary: '#605EA1',
          secondary: '#22177A',
          neutral: '#8EA3A6', 
          accent: '#E6E9AF',
        },
        two: {
          primary: '#FFEBB7',
          secondary: '#243763',
          neutral: '#AD8E70', 
          accent: '#FF6E31',
        },
        three: {
          primary: '#FFE6E6',
          secondary: '#E1AFD1',
          neutral: '#AD88C6', 
          accent: '#7469B6',
        },
        four: {
          primary: '#F5EEE6',
          secondary: '#FFF8E3',
          neutral: '#F3D7CA', 
          accent: '#E6A4B4',
        },
        five: {
          primary: '#DCF2F1',
          secondary: '#7FC7D9',
          neutral: '#365486', 
          accent: '#365486',
        },
        six: {
          primary: '#F9F7F7',
          secondary: '#DBE2EF',
          neutral: '#3F72AF', 
          accent: '#112D4E',
        },
        seven: {
          primary: '#FEFAF6',
          secondary: '#EADBC8',
          neutral: '#DAC0A3', 
          accent: '#102C57',
        }
      },
    ]
  },
  theme: {
    extend: {
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

