Type,Subtype,Feature
Stage,General,add_fields
Stage,General,aggregate
Stage,General,collection
Stage,General,collection_group
Stage,General,database
Stage,General,distinct
Stage,General,documents
Stage,General,find_nearest
Stage,General,limit
Stage,General,literals
Stage,General,offset
Stage,General,remove_fields
Stage,General,select
Stage,General,sort
Stage,General,where
Stage,General,replace_with
Stage,General,sample 
Stage,General,union
Stage,General,unnest
Stage,General,let
Stage,General,subcollection
Stage,General,search
Stage,General,count
Stage,General,set_default
Stage,General,sample (system)
Stage,General,stats
Stage,General,window
Stage,General,paginate
Stage,General,lookup
Stage,General,bucket
Stage,General,limit_to_last
Stage,General,graph / recursive CTEs
Stage,DML,update
Stage,DML,upsert
Stage,DML,delete
Stage,DML,insert
Function,General,get_field
Function,General,offset
Function,Accumulators (Aggregation),average
Function,Accumulators (Aggregation),count
Function,Accumulators (Aggregation),count_if
Function,Accumulators (Aggregation),maximum
Function,Accumulators (Aggregation),minimum
Function,Accumulators (Aggregation),sum
Function,Accumulators (Aggregation),count_distinct
Function,Accumulators (Aggregation),first
Function,Accumulators (Aggregation),last
Function,Accumulators (Aggregation),array_agg
Function,Accumulators (Aggregation),array_agg_distinct
Function,Accumulators (Aggregation),first_n
Function,Accumulators (Aggregation),last_n
Function,Accumulators (Aggregation),maximum_n
Function,Accumulators (Aggregation),minimum_n
Function,Accumulators (Aggregation),approx_count_distinct
Function,Accumulators (Aggregation),median
Function,Accumulators (Aggregation),percentile
Function,Accumulators (Aggregation),std_dev_pop
Function,Accumulators (Aggregation),std_dev_samp
Function,Accumulators (Aggregation),logical_or
Function,Accumulators (Aggregation),logical_and
Function,Accumulators (Aggregation),any
Function,Accumulators (Aggregation),bottom
Function,Accumulators (Aggregation),bottom_n
Function,Accumulators (Aggregation),top
Function,Accumulators (Aggregation),top_n
Function,Arithmetic,add
Function,Arithmetic,divide
Function,Arithmetic,multiply
Function,Arithmetic,subtract
Function,Arithmetic,abs
Function,Arithmetic,ceil
Function,Arithmetic,exp
Function,Arithmetic,floor
Function,Arithmetic,ln
Function,Arithmetic,log
Function,Arithmetic,log10
Function,Arithmetic,mod
Function,Arithmetic,pow
Function,Arithmetic,round
Function,Arithmetic,sqrt
Function,Arithmetic,rand
Function,Arithmetic,trunc
Function,Array,array_contains
Function,Array,array_contains_all
Function,Array,array_contains_any
Function,Array,array_length
Function,Array,equal_any
Function,Array,not_equal_any
Function,Array,array
Function,Array,array_get
Function,Array,array_reverse
Function,Array,array_concat
Function,Array,sum
Function,Array,maximum
Function,Array,minimum
Function,Array,first
Function,Array,first_n
Function,Array,last
Function,Array,last_n
Function,Array,maximum_n
Function,Array,minimum_n
Function,Array,array_index_of
Function,Array,array_index_of_all
Function,Array,array_slice
Function,Array,array_filter
Function,Array,array_transform
Function,Array,avg
Function,Array,array_find
Function,Array,array_sort
Function,Array,array_flatten
Function,Array,array_is_distinct
Function,Array,range
Function,Array,array_reduce
Function,Array,array_zip
Function,Bitwise,bit_and
Function,Bitwise,bit_not
Function,Bitwise,bit_or
Function,Bitwise,bit_xor
Function,Comparison,equal
Function,Comparison,greater_than
Function,Comparison,greater_than_or_equal
Function,Comparison,less_than
Function,Comparison,less_than_or_equal
Function,Comparison,not_equal
Function,Comparison,cmp
Function,Data Size,storage_size
Function,Data Size,document_size
Function,Date / Timestamp,timestamp_add
Function,Date / Timestamp,timestamp_subtract
Function,Date / Timestamp,timestamp_to_unix_micros
Function,Date / Timestamp,timestamp_to_unix_millis
Function,Date / Timestamp,timestamp_to_unix_seconds
Function,Date / Timestamp,unix_micros_to_timestamp
Function,Date / Timestamp,unix_millis_to_timestamp
Function,Date / Timestamp,unix_seconds_to_timestamp
Function,Date / Timestamp,current_timestamp
Function,Date / Timestamp,timestamp_trunc
Function,Date / Timestamp,timestamp_diff
Function,Date / Timestamp,timestamp_extract
Function,Date / Timestamp,format_timestamp
Function,Date / Timestamp,parse_timestamp
Function,General,length
Function,General,reverse
Function,General,concat
Function,General,current_document
Function,Geospatial,geo_distance
Function,Key,collection_id
Function,Key,document_id
Function,Key,namespace
Function,Key,parent
Function,Key,has_ancestor
Function,Key,reference_slice
Function,Key,replace_document_id
Function,Key,autogen_document_id
Function,Logical,and
Function,Logical,exists
Function,Logical,not
Function,Logical,or
Function,Logical,xor
Function,Logical,nor
Function,Logical,conditional
Function,Logical,if_null
Function,Logical,coalesce
Function,Logical,logical_max
Function,Logical,logical_min
Function,Logical,is_error
Function,Logical,if_error
Function,Logical,if_absent
Function,Logical,is_absent
Function,Logical,error
Function,Logical,switch_on
Function,Object,map
Function,Object,map_get
Function,Object,map_set
Function,Object,map_merge
Function,Object,map_remove
Function,Object,map_keys
Function,Object,map_values
Function,Object,map_entries
Function,Object,map_mask
Function,Set,all_elements_true
Function,Set,any_elements_true
Function,Set,set_difference
Function,Set,set_equals
Function,Set,set_intersection
Function,Set,set_is_subset
Function,Set,set_union
Function,String,byte_length
Function,String,char_length
Function,String,ends_with
Function,String,like
Function,String,regex_contains
Function,String,regex_match
Function,String,regex_find
Function,String,regex_find_all
Function,String,starts_with
Function,String,string_concat
Function,String,string_contains
Function,String,string_reverse
Function,String,join
Function,String,substring
Function,String,to_lower
Function,String,to_upper
Function,String,trim
Function,String,string_split
Function,String,regex_replace
Function,String,regex_replace_first
Function,String,string_repeat
Function,String,string_replace_all
Function,String,string_replace_one
Function,String,string_index_of
Function,String,ltrim
Function,String,rtrim
Function,String,string_case_cmp
Function,Trigonometry,acos
Function,Trigonometry,acosh
Function,Trigonometry,asin
Function,Trigonometry,asinh
Function,Trigonometry,atan
Function,Trigonometry,atan2
Function,Trigonometry,atanh
Function,Trigonometry,cos
Function,Trigonometry,cosh
Function,Trigonometry,degrees_to_radians
Function,Trigonometry,radians_to_degrees
Function,Trigonometry,sin
Function,Trigonometry,sinh
Function,Trigonometry,tan
Function,Trigonometry,tanh
Function,Type,type
Function,Type,is_type
Function,Type,cast
Function,Vector,cosine_distance
Function,Vector,dot_product
Function,Vector,euclidean_distance
Function,Vector,vector_length
Function,General,let