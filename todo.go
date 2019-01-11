package main

// Unconverted stuff that we may need...



// this function determines OPTIMAL_PROTO. It is a workaround function as under certain circumstances
// (e.g. IIS6.0 and openssl 1.0.2 as opposed to 1.0.1) needs a protocol otherwise s_client -connect will fail!
// Circumstances observed so far: 1.) IIS 6  2.) starttls + dovecot imap
// The first try in the loop is empty as we prefer not to specify always a protocol if we can get along w/o it
//
func determineOptimalProtocol(proto string) string {
/*
local all_failed=true
local tmp=""

>$ERRFILE
if [[ -n "$1" ]]; then
// starttls workaround needed see https://github.com/drwetter/testssl.sh/issues/188 -- kind of odd
for STARTTLS_OPTIMAL_PROTO in -tls1_2 -tls1 -ssl3 -tls1_1 -tls1_3 -ssl2; do
case $STARTTLS_OPTIMAL_PROTO in
-tls1_3) "$HAS_TLS13" || continue ;;
-ssl3)   "$HAS_SSL3" || continue ;;
-ssl2)   "$HAS_SSL2" || continue ;;
*) ;;
esac
$OPENSSL s_client $(s_client_options "$STARTTLS_OPTIMAL_PROTO $BUGS -connect "$NODEIP:$PORT" $PROXY -msg -starttls $1") </dev/null >$TMPFILE 2>>$ERRFILE
if sclient_auth $? $TMPFILE; then
all_failed=false
break
fi
all_failed=true
done
"$all_failed" && STARTTLS_OPTIMAL_PROTO=""
debugme echo "STARTTLS_OPTIMAL_PROTO: $STARTTLS_OPTIMAL_PROTO"
else
for OPTIMAL_PROTO in '' -tls1_2 -tls1 -tls1_3 -ssl3 -tls1_1 -ssl2; do
case $OPTIMAL_PROTO in
-tls1_3) "$HAS_TLS13" || continue ;;
-ssl3)   "$HAS_SSL3" || continue ;;
-ssl2)   "$HAS_SSL2" || continue ;;
*) ;;
esac
$OPENSSL s_client $(s_client_options "$OPTIMAL_PROTO $BUGS -connect "$NODEIP:$PORT" -msg $PROXY $SNI") </dev/null >$TMPFILE 2>>$ERRFILE
if sclient_auth $? $TMPFILE; then
// we use the successful handshake at least to get one valid protocol supported -- it saves us time later
if [[ -z "$OPTIMAL_PROTO" ]]; then
// convert to openssl terminology
tmp=$(get_protocol $TMPFILE)
tmp=${tmp/\./_}
tmp=${tmp/v/}
tmp="$(tolower $tmp)"
add_tls_offered "${tmp}" yes
else
add_tls_offered "${OPTIMAL_PROTO/-/}" yes
fi
debugme echo "one proto determined: $tmp"
all_failed=false
break
fi
all_failed=true
done
"$all_failed" && OPTIMAL_PROTO=""
debugme echo "OPTIMAL_PROTO: $OPTIMAL_PROTO"
if [[ "$OPTIMAL_PROTO" == "-ssl2" ]]; then
prln_magenta "$NODEIP:$PORT appears to only support SSLv2."
ignore_no_or_lame " Type \"yes\" to proceed and accept false negatives or positives" "yes"
[[ $? -ne 0 ]] && exit $ERR_CLUELESS
fi
fi
grep -q '^Server Temp Key' $TMPFILE && HAS_DH_BITS=true     // FIX //190

if "$all_failed"; then
outln
if "$HAS_IPv6"; then
pr_bold " Your $OPENSSL is not IPv6 aware, or $NODEIP:$PORT "
else
pr_bold " $NODEIP:$PORT "
fi
tmpfile_handle ${FUNCNAME[0]}.txt
prln_bold "doesn't seem to be a TLS/SSL enabled server";
ignore_no_or_lame " The results might look ok but they could be nonsense. Really proceed ? (\"yes\" to continue)" "yes"
[[ $? -ne 0 ]] && exit $ERR_CLUELESS
fi

// NOTE: The following code is only needed as long as draft versions of TLSv1.3 prior to draft 23
// are supported. It is used to determine whether a draft 23 or pre-draft 23 ClientHello should be
// sent.
if [[ -z "$1" ]]; then
KEY_SHARE_EXTN_NR="33"
tls_sockets "04" "$TLS13_CIPHER" "" "00, 2b, 00, 0f, 0e, 03,04, 7f,1c, 7f,1b, 7f,1a, 7f,19, 7f,18, 7f,17"
if [[ $? -eq 0 ]]; then
add_tls_offered tls1_3 yes
else
KEY_SHARE_EXTN_NR="28"
tls_sockets "04" "$TLS13_CIPHER" "" "00, 2b, 00, 0b, 0a, 7f,16, 7f,15, 7f,14, 7f,13, 7f,12"
if [[ $? -eq 0 ]]; then
add_tls_offered tls1_3 yes
else
add_tls_offered tls1_3 no
KEY_SHARE_EXTN_NR="33"
fi
fi
fi

tmpfile_handle ${FUNCNAME[0]}.txt
*/
return proto
}



