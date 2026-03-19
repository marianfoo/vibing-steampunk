CLASS zcl_lexer DEFINITION PUBLIC FINAL CREATE PUBLIC.
  PUBLIC SECTION.
    METHODS run IMPORTING iv_source TYPE string RETURNING VALUE(rv_result) TYPE TABLE OF REF TO zcl_token.
  PRIVATE SECTION.
    DATA mt_tokens TYPE TABLE OF REF TO zcl_token.
ENDCLASS.

CLASS zcl_lexer IMPLEMENTATION.
  METHOD run.
    me->mv_tokens = VALUE #( ).
    DATA(lv_pos) = 0.
    DATA(lv_len) = strlen( lv_source ).
    WHILE lv_pos < lv_len.
      DATA(lv_ch) = lv_source+lv_pos(1).
      IF lv_ch = ' ' OR lv_ch = '
' OR lv_ch = '	'.
        lv_pos = lv_pos + 1.
        CONTINUE.
      ENDIF.
      IF lv_ch = '*' AND lv_pos = 0.
        DATA(lv_end) = find( val = lv_source sub = '
' off = lv_pos ).
        IF lv_end = - 1.
          lv_end = lv_len.
        ENDIF.
        DATA(lv_comment) = substring( val = lv_source off = lv_pos len = lv_end - lv_pos ).
        APPEND NEW zcl_token( iv_p1 = NEW zcl_position( iv_p1 = 1 iv_p2 = lv_pos ) iv_p2 = lv_comment ) TO me->mv_tokens.
        lv_pos = lv_end.
        CONTINUE.
      ENDIF.
      IF lv_ch = '"'.
        DATA(lv_end) = find( val = lv_source sub = '
' off = lv_pos ).
        IF lv_end = - 1.
          lv_end = lv_len.
        ENDIF.
        DATA(lv_comment) = substring( val = lv_source off = lv_pos len = lv_end - lv_pos ).
        APPEND NEW zcl_token( iv_p1 = NEW zcl_position( iv_p1 = 1 iv_p2 = lv_pos ) iv_p2 = lv_comment ) TO me->mv_tokens.
        lv_pos = lv_end.
        CONTINUE.
      ENDIF.
      IF lv_ch = ''''.
        DATA(lv_end) = find( val = lv_source sub = '''' off = lv_pos + 1 ).
        IF lv_end = - 1.
          lv_end = lv_len.
        ENDIF.
        DATA(lv_str) = substring( val = lv_source off = lv_pos len = lv_end + 1 - lv_pos ).
        APPEND NEW zcl_token( iv_p1 = NEW zcl_position( iv_p1 = 1 iv_p2 = lv_pos ) iv_p2 = lv_str ) TO me->mv_tokens.
        lv_pos = lv_end + 1.
        CONTINUE.
      ENDIF.
      IF lv_ch = '.' OR lv_ch = ',' OR lv_ch = ':' OR lv_ch = '(' OR lv_ch = ')' OR lv_ch = '[' OR lv_ch = ']'.
        APPEND NEW zcl_token( iv_p1 = NEW zcl_position( iv_p1 = 1 iv_p2 = lv_pos ) iv_p2 = lv_ch ) TO me->mv_tokens.
        lv_pos = lv_pos + 1.
        CONTINUE.
      ENDIF.
      DATA(lv_end) = lv_pos + 1.
      WHILE lv_end < lv_len.
        DATA(lv_c) = lv_source+lv_end(1).
        IF lv_c = ' ' OR lv_c = '
' OR lv_c = '	' OR lv_c = '.' OR lv_c = ',' OR lv_c = ':' OR lv_c = '(' OR lv_c = ')' OR lv_c = '[' OR lv_c = ']'.
          EXIT.
        ENDIF.
        lv_end = lv_end + 1.
      ENDWHILE.
      APPEND NEW zcl_token( iv_p1 = NEW zcl_position( iv_p1 = 1 iv_p2 = lv_pos ) iv_p2 = substring( val = lv_source off = lv_pos len = lv_end - lv_pos ) ) TO me->mv_tokens.
      lv_pos = lv_end.
    ENDWHILE.
    rv_result = me->mv_tokens.
    RETURN.
  ENDMETHOD.
ENDCLASS.
