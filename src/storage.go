package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sqlx.DB
}

var schema = `
CREATE TABLE IF NOT EXISTS public.globus_item (
	created_at timestamp,
    item_id text,
    item_code text,
    product_id text,
	item_name text,
	href text,
	picture text,
	price text,
	old_price text,
	step text,
	item_min text,
	item_max text,
	item_type text,
	measure text,
	measure_ratio text,
	show_discount bool,
	is_loyal_price bool,
	show_red_price bool,
	discount_percent text,
	av_weight text,
	is_price_gram bool,
	is_alcogol bool,
	is_alco_allowed bool,
	is_signable bool,
	chips_price text,
	item_list text,
	brand text,
	category text,
	is_theme bool,
	sort_section text,
	print_price text,
	print_old_price text,
	xml_id text,
	hyper_code text,
	is_marked bool,
	section_id text,
	section_name text
);
`

func (s *Store) Open() error {
	db, err := sqlx.Open("postgres", "host=apps dbname=parser sslmode=disable user=postgres password=<pass>")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db

	db.MustExec(schema)

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) SaveGlobus(items map[string]GlobusItem) {

	list := make([]GlobusItem, 0, len(items))
	for _, v := range items {
		list = append(list, v)
	}
	_, err := s.db.NamedExec(`INSERT INTO public.globus_item (
			created_at, item_id, item_code, product_id, item_name, href, picture, price, old_price, step, item_min, item_max, item_type, measure,
			measure_ratio,
			show_discount,
			is_loyal_price,
			show_red_price,
			discount_percent,
			av_weight,
			is_price_gram,
			is_alcogol,
			is_alco_allowed,
			is_signable,
			chips_price,
			item_list,
			brand,
			category,
			is_theme,
			sort_section,
			print_price,
			print_old_price,
			xml_id,
			hyper_code,
			is_marked,
			section_id,
			section_name
		) VALUES (
			now(),
			:item_id,
			:item_code,
			:product_id,
			:item_name,
			:href,
			:picture,
			:price,
			:old_price,
			:step,
			:item_min,
			:item_max,
			:item_type,
			:measure,
			:measure_ratio,
			:show_discount,
			:is_loyal_price,
			:show_red_price,
			:discount_percent,
			:av_weight,
			:is_price_gram,
			:is_alcogol,
			:is_alco_allowed,
			:is_signable,
			:chips_price,
			:item_list,
			:brand,
			:category,
			:is_theme,
			:sort_section,
			:print_price,
			:print_old_price,
			:xml_id,
			:hyper_code,
			:is_marked,
			:section_id,
			:section_name
		)`, list)

	if err != nil {
		log.Panic(err)
	}
}
