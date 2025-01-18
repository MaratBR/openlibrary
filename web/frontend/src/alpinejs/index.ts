import Alpine from 'alpinejs'
//@ts-expect-error
import ajax from '@imacrayon/alpine-ajax'

// alpinejs components
import './book-reader'
import './rating-input'


(window as any).Alpine = Alpine;
Alpine.plugin(ajax)
Alpine.start();
